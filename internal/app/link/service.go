package link

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"github.com/teris-io/shortid"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Add(*gin.Context, addRequest) (model.Link, error)
	FindByID(*gin.Context, findByIDRequest) (model.Link, error)
	FindAll(*gin.Context, findAllRequest) (int64, int, []model.Link, error)
	Update(*gin.Context, updateRequest) error
	Delete(*gin.Context, deleteRequest) (err error)
	FindFullURL(*gin.Context, string, string) (model.Link, error)

	NewHTTP(*gin.RouterGroup)
	HTTPAdd(*gin.Context)
	HTTPFindByID(*gin.Context)
	HTTPFindAll(*gin.Context)
	HTTPUpdate(*gin.Context)
	HTTPDelete(*gin.Context)
	HTTPFindFullURL(*gin.Context)
}

type service struct {
	db    *sql.DB
	cache *redis.Client
}

// NewService creates a new authentication service.
func NewService(db *sql.DB, cache *redis.Client) Service {
	return service{db, cache}
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
		ORDER BY id ASC OFFSET $2 LIMIT $3
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

		links = append(links, link)
	}

	fmt.Println("LINKS: ", links)

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

	key := fmt.Sprintf("link_full_%s_%s", domain, keyword)
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

	log.Debug().Caller().Msg("No, key '%s' is not cached!")
	return
}
