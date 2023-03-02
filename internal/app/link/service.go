package link

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/wvoliveira/corgi/internal/pkg/common"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"github.com/teris-io/shortid"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

const (
	keyCacheShortLink                   = "cache:link_short:%s:%s"                      // Ex.: cache:link_short:domain:keyword
	keyCacheShortLinkMetricCounterTotal = "cache:link_short:%s:%s:metric:counter_total" // Cache value for sum of counter metric.

	// These values stay inside hash value.
	keyMetricShortLink = "metric:link_short:%s:%s" // Ex.: metric:link_short:domain:keyword
	keyMetricCounter   = "counter:%s"              // Ex.: counter:yyyy-mm-dd-hh-MM
	keyLatestSync      = "latest:sync:%s"          // Ex.: latest:sync:<timestamp>
	keyLatestCheck     = "latest:check:%s"         // Ex.: latest:check:<timestamp>
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	FindRedirectURL(*gin.Context, string, string) (model.Link, error)
	Add(*gin.Context, addRequest) (model.Link, error)
	FindByID(*gin.Context, findByIDRequest) (model.Link, error)
	FindAll(*gin.Context, findAllRequest) (int64, int, []model.Link, error)
	Update(*gin.Context, updateRequest) error
	Delete(*gin.Context, deleteRequest) (err error)
	FindFullURL(*gin.Context, string, string) (model.Link, error)
	Clicks(*gin.Context, clicksRequest) (model.LinkClicks, error)

	NewHTTP(*gin.Engine, *gin.RouterGroup)
	HTTPRedirect(*gin.Context)
	HTTPAdd(*gin.Context)
	HTTPFindByID(*gin.Context)
	HTTPFindAll(*gin.Context)
	HTTPUpdate(*gin.Context)
	HTTPDelete(*gin.Context)
	HTTPFindFullURL(*gin.Context)
	HTTPClicks(*gin.Context)
}

type service struct {
	db    *sql.DB
	cache *redis.Client
}

// NewService creates a new authentication service.
func NewService(db *sql.DB, cache *redis.Client) Service {
	return service{db, cache}
}

// FindRedirectURL redirect to full link getting by domain and keyword combination.
func (s service) FindRedirectURL(ctx *gin.Context, domain, keyword string) (m model.Link, err error) {
	log := logger.Logger(ctx)

	key := fmt.Sprintf(keyCacheShortLink, domain, keyword)
	val, _ := itemFromCache(ctx, s.cache, key)
	if val != "" {
		m.URL = val

		// If found, increase counter in background process.
		go increaseCounter(ctx, s.cache, domain, keyword)
		return
	}

	query := "SELECT url FROM links WHERE domain = $1 AND keyword = $2 AND active = true"
	log.Debug().Caller().Msg(query)

	err = s.db.QueryRowContext(ctx, query, domain, keyword).Scan(&m.URL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return m, e.ErrLinkNotFound
		}

		log.Error().Caller().Msg(err.Error())
		return
	}

	// If found, increase counter in background process.
	go increaseCounter(ctx, s.cache, domain, keyword)

	status := s.cache.Set(ctx, key, m.URL, 10*time.Minute)
	err = status.Err()
	if err != nil {
		log.Error().Caller().Msg(err.Error())
	}
	return
}

// Add create a new shortener link.
func (s service) Add(c *gin.Context, payload addRequest) (link model.Link, err error) {
	log := logger.Logger(c)

	if err = checkLink(payload.Domain, payload.Keyword, payload.URL); err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	// If user is anonymous, create a random ID and blank another fields.
	if payload.WhoID == "0" {
		sid, _ := shortid.New(1, shortid.DefaultABC, 2342)
		payload.Keyword, _ = sid.Generate()
	}

	// If not anonymous access, create a random keyword if was not set.
	if payload.WhoID != "0" {
		if payload.Keyword == "" {
			sid, _ := shortid.New(1, shortid.DefaultABC, 2342)
			payload.Keyword, _ = sid.Generate()
		}
	}

	query := "SELECT id FROM links WHERE domain = $1 AND keyword = $2 LIMIT 1"
	err = s.db.QueryRowContext(c, query, payload.Domain, payload.Keyword).Scan(&link.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Error().Caller().Msg(err.Error())
			return
		}
	}

	if link.ID != "" {
		message := fmt.Sprintf("link with domain '%s' and keyword '%s' already exists", payload.Domain, payload.Keyword)
		log.Warn().Caller().Msg(message)
		return link, e.ErrLinkAlreadyExists
	}

	query = `
		INSERT INTO links(id, domain, keyword, url, title, user_id) 
		VALUES($1, $2, $3, $4, $5, $6)
	`

	// Create a new link getting info from payload.
	// Maybe we can change this to a more elegant way.
	newLink := model.Link{}
	newLink.ID = ulid.Make().String()
	newLink.Domain = payload.Domain
	newLink.Keyword = payload.Keyword
	newLink.URL = payload.URL
	newLink.Title = payload.Title
	newLink.UserID = payload.WhoID

	_, err = s.db.ExecContext(
		c,
		query,
		newLink.ID,
		newLink.Domain,
		newLink.Keyword,
		newLink.URL,
		newLink.Title,
		newLink.UserID,
	)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	query = `SELECT id, user_id, created_at, updated_at, domain, keyword, url, title, active FROM links WHERE id = $1`
	err = s.db.QueryRowContext(c, query, newLink.ID).Scan(
		&link.ID,
		&link.UserID,
		&link.CreatedAt,
		&link.UpdatedAtNull,
		&link.Domain,
		&link.Keyword,
		&link.URL,
		&link.Title,
		&link.Active)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}
	return
}

// FindByID get a shortener link from ID.
func (s service) FindByID(c *gin.Context, payload findByIDRequest) (link model.Link, err error) {
	log := logger.Logger(c)

	query := "SELECT * FROM links WHERE user_id = $1 AND id = $2 LIMIT 1"
	rows, err := s.db.QueryContext(c, query, payload.WhoID, payload.LinkID)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(
			&link.ID,
			&link.UserID,
			&link.CreatedAt,
			&link.UpdatedAt,
			&link.Domain,
			&link.Keyword,
			&link.URL,
			&link.Title,
			&link.Active,
		)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return
		}

		log.Debug().Caller().Msg(fmt.Sprintf("link_id=%s", link.ID))
		return
	}

	log.Debug().Caller().Msg("link not found")
	return link, e.ErrLinkNotFound
}

// FindAll get a list of links from database.
func (s service) FindAll(ctx *gin.Context, payload findAllRequest) (total int64, pages int, links []model.Link, err error) {
	log := logger.Logger(ctx)

	queryCount := `SELECT COUNT(0) FROM links 
                WHERE user_id = $1
	`
	log.Debug().Caller().Msg(queryCount)

	queryData := `SELECT id, user_id, created_at, updated_at, domain, keyword, url, title, active 
		FROM links
		WHERE user_id = $1
		ORDER BY $2 OFFSET $3 LIMIT $4
	`
	log.Debug().Caller().Msg(queryData)

	//domain, keyword := common.SplitURL(payload.ShortenedURL)

	err = s.db.QueryRowContext(
		ctx,
		queryCount,
		payload.WhoID,
		//payload.SearchText,
		//domain,
		//keyword,
	).Scan(&total)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	rows, err := s.db.QueryContext(
		ctx,
		queryData,
		payload.WhoID,
		payload.Sort,
		payload.Offset,
		payload.Limit,
	)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	defer rows.Close()
	links = []model.Link{}
	link := model.Link{}

	for rows.Next() {
		err = rows.Scan(
			&link.ID,
			&link.UserID,
			&link.CreatedAt,
			&link.UpdatedAtNull,
			&link.Domain,
			&link.Keyword,
			&link.URL,
			&link.Title,
			&link.Active,
		)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return
		}

		if link.UpdatedAtNull.Valid {
			link.UpdatedAt = &link.UpdatedAtNull.Time
		}

		//clicks, err := s.Clicks(ctx, clicksRequest{
		//	WhoID:    payload.WhoID,
		//	ShortURL: fmt.Sprintf("%s/%s", link.Domain, link.Keyword),
		//})
		//if err == nil {
		//	link.Clicks = clicks
		//}

		links = append(links, link)
	}

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return total, pages, links, e.ErrInternalServerError
	}

	pages = int(math.Ceil(float64(total) / float64(payload.Limit)))
	return
}

// Update change specific link by ID.
func (s service) Update(ctx *gin.Context, payload updateRequest) (err error) {
	log := logger.Logger(ctx)

	query := "UPDATE links SET title = $1 WHERE id = $2 AND user_id = $3"
	log.Debug().Caller().Msg(query)

	_, err = s.db.ExecContext(ctx, query, payload.Title, payload.LinkID, payload.WhoID)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return e.ErrInternalServerError
	}
	return
}

// Delete delete a link by ID.
func (s service) Delete(ctx *gin.Context, payload deleteRequest) (err error) {
	log := logger.Logger(ctx)
	link := model.Link{}

	query := "SELECT id, domain, keyword FROM links WHERE user_id = $1 AND id = $2 AND active = true"
	err = s.db.QueryRowContext(ctx, query, payload.WhoID, payload.LinkID).Scan(&link.ID, &link.Domain, &link.Keyword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug().Caller().Msg(fmt.Sprintf("Link with id=%s was not found", payload.LinkID))
			return e.ErrLinkNotFound
		}

		log.Error().Caller().Msg(err.Error())
		return e.ErrInternalServerError
	}

	key := fmt.Sprintf("link_%s_%s", link.Domain, link.Keyword)
	_, err = s.cache.Del(ctx, key).Result()

	// Keep going on error from cache.
	// Because SQL database still working, so cache doesn't matter at this moment.
	if err != nil {
		log.Error().Caller().Msg(err.Error())
	}

	query = "UPDATE links SET active = false, updated_at = $1 WHERE user_id = $2 AND id = $3"
	log.Debug().Caller().Msg(query)

	_, err = s.db.ExecContext(ctx, query, time.Now(), payload.WhoID, payload.LinkID)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return e.ErrInternalServerError
	}
	return
}

// FindFullURL get a shortener link from keyword.
func (s service) FindFullURL(c *gin.Context, domain, keyword string) (m model.Link, err error) {
	log := logger.Logger(c)

	key := fmt.Sprintf(keyCacheShortLink, domain, keyword)
	val, _ := itemFromCache(c, s.cache, key)
	if val != "" {
		m.URL = val
		return
	}

	query := "SELECT url FROM links WHERE domain = $1 AND keyword = $2 AND active = true"
	log.Debug().Caller().Msg(query)

	err = s.db.QueryRowContext(c, query, domain, keyword).Scan(&m.URL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return m, e.ErrLinkNotFound
		}

		log.Error().Caller().Msg(err.Error())
		return
	}

	status := s.cache.Set(c, key, m.URL, 10*time.Minute)
	err = status.Err()
	if err != nil {
		fmt.Println(err.Error())
	}
	return
}

func (s service) Clicks(ctx *gin.Context, payload clicksRequest) (lc model.LinkClicks, err error) {
	log := logger.Logger(ctx)

	// Check if this key exists in cache,
	// if not, set with expiration for 10 seconds.
	domain, keyword := common.SplitURL(payload.ShortURL)
	keyCache := fmt.Sprintf(keyCacheShortLinkMetricCounterTotal, domain, keyword)

	val, _ := itemFromCache(ctx, s.cache, keyCache)
	if val != "" {
		total, err := strconv.Atoi(val)
		if err != nil {
			return lc, err
		}

		lc.Total = total
		return lc, err
	}

	key := fmt.Sprintf(keyMetricShortLink, domain, keyword)
	log.Debug().Caller().Msg(key)

	cmd := s.cache.HVals(ctx, key)
	result, err := cmd.Result()
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	total := 0
	for _, value := range result {
		i, _ := strconv.Atoi(value)
		total += i
	}

	// Send item to cache async.
	// But whatever if errors happens.
	go func() {
		err := setItemToCache(ctx, s.cache, keyCache, total)
		if err != nil {
			log.Error().Caller().Msg(fmt.Sprintf("Error to create cache: %s", err.Error()))
		}
	}()

	lc.Total = total
	return
}

// itemFromCache get value from cache with specific key.
func itemFromCache(c context.Context, cache *redis.Client, key string) (item string, err error) {
	log := logger.Logger(c)
	log.Debug().Caller().Msg(fmt.Sprintf("Collecting key '%s' from cache", key))

	item, err = cache.Get(c, key).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Error().Caller().Msg(err.Error())
			return
		}
	}

	if item != "" {
		log.Debug().Caller().Msg(fmt.Sprintf("OK, I gotten key '%s' from cache!", key))
		return
	}

	log.Debug().Caller().Msg(fmt.Sprintf("Key '%s' is not cached!", key))
	return
}

// setItemToCache send item to cache and set 10 seconds for default expiration.
func setItemToCache(ctx context.Context, cache *redis.Client, key string, value any) (err error) {
	status := cache.Set(ctx, key, value, time.Second*10)
	err = status.Err()
	if err != nil {
		return
	}
	return
}

// increaseCounter sum by 1 with Redis hash type.
func increaseCounter(ctx context.Context, cache *redis.Client, domain, keyword string) {
	log := logger.Logger(ctx)

	now := time.Now().Format("2006-01-02-15")
	log.Debug().Caller().Msg(fmt.Sprintf("Now: %s", now))

	keyHash := fmt.Sprintf(keyMetricShortLink, domain, keyword)
	keyCounter := fmt.Sprintf(keyMetricCounter, now)

	log.Debug().Caller().Msg(fmt.Sprintf("Key hash: %s", keyHash))
	log.Debug().Caller().Msg(fmt.Sprintf("Key counter: %s", keyCounter))

	cache.HSet(ctx, keyHash)
	stat := cache.HIncrBy(ctx, keyHash, keyCounter, 1)

	counter, _ := stat.Result()
	log.Debug().Caller().Msg(fmt.Sprintf("Counter per hour: %d", counter))
}
