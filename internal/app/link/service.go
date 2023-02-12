package link

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/teris-io/shortid"
	"github.com/wvoliveira/corgi/internal/pkg/common"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Add(*gin.Context, model.Link) (model.Link, error)
	FindByID(*gin.Context, string, string) (model.Link, error)
	FindAll(*gin.Context, findAllRequest) (int64, int, []model.Link, error)
	Update(*gin.Context, model.Link) error
	Delete(*gin.Context, string, string) (err error)

	NewHTTP(*gin.RouterGroup)
	HTTPAdd(*gin.Context)
	HTTPFindByID(*gin.Context)
	HTTPFindAll(*gin.Context)
	HTTPUpdate(*gin.Context)
	HTTPDelete(*gin.Context)
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
func (s service) Add(c *gin.Context, payload model.Link) (m model.Link, err error) {
	l := logger.Logger(c)

	if err = checkLink(payload); err != nil {
		l.Error().Caller().Msg(err.Error())
		return
	}

	// If user is anonymous, create a random ID and blank another fields.
	if payload.UserID == "0" {
		sid, _ := shortid.New(1, shortid.DefaultABC, 2342)
		payload.Keyword, _ = sid.Generate()
	}

	if payload.UserID != "0" {
		if payload.Keyword == "" {
			sid, _ := shortid.New(1, shortid.DefaultABC, 2342)
			payload.Keyword, _ = sid.Generate()
		}
	}

	query := "SELECT id FROM links WHERE domain = $1 AND keyword = $2 LIMIT 1"
	err = s.db.QueryRowContext(c, query, payload.Domain, payload.Keyword).Scan(&m.ID)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Error().Caller().Msg(err.Error())
			return
		}
	}

	if m.ID != "" {
		l.Warn().Caller().Msg("domain with keyword already exists")
		return m, e.ErrLinkAlreadyExists
	}

	query = `
		INSERT INTO links(id, domain, keyword, url, title, user_id) 
		VALUES($1, $2, $3, $4, $5, $6)
	`

	payload.ID = ulid.Make().String()

	_, err = s.db.ExecContext(
		c,
		query,
		payload.ID,
		payload.Domain,
		payload.Keyword,
		payload.URL,
		payload.Title,
		payload.UserID,
	)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	err = s.db.QueryRowContext(c, "SELECT * FROM links WHERE id = $1", payload.ID).Scan(
		&m.ID, &m.UserID, &m.CreatedAt, &m.UpdatedAt, &m.Domain, &m.Keyword, &m.URL, &m.Title, &m.Active)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	return
}

// FindByID get a shortener link from ID.
func (s service) FindByID(c *gin.Context, linkID, userID string) (link model.Link, err error) {
	log := logger.Logger(c)

	query := "SELECT * FROM links WHERE id = $1 AND user_id = $2 LIMIT 1"
	rows, err := s.db.QueryContext(c, query, linkID, userID)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&link.ID, &link.UserID, &link.CreatedAt, &link.UpdatedAt, &link.Domain, &link.Keyword, &link.URL, &link.Title, &link.Active)

		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return
		}

		log.Info().Caller().Msg(fmt.Sprintf("link ID: %s", link.ID))
		return
	}

	log.Warn().Caller().Msg("link not found")
	return link, e.ErrLinkNotFound
}

// FindAll get a list of links from database.
func (s service) FindAll(c *gin.Context, r findAllRequest) (total int64, pages int, links []model.Link, err error) {
	log := logger.Logger(c)

	queryCount := "SELECT COUNT(0) FROM links "
	queryData := "SELECT * FROM links "

	queryFilter := fmt.Sprintf(" WHERE user_id = '%s'", r.UserID)

	if len(r.SearchText) >= 3 {
		queryFilter = queryFilter + fmt.Sprintf(" AND domain LIKE '%%%[1]s%%' OR keyword LIKE '%%%[1]s%%' ", r.SearchText)
	}

	domain, keyword := common.SplitURL(r.ShortenedURL)
	if domain != "" && keyword != "" {
		queryFilter = queryFilter + fmt.Sprintf(" AND domain = '%s' AND keyword = '%s' ", domain, keyword)
	}

	// TODO: add order by another field.
	queryCount = queryCount + queryFilter
	queryData = queryData + queryFilter + fmt.Sprintf(" ORDER BY ID DESC OFFSET %d LIMIT %d ", r.Offset, r.Limit)

	err = s.db.QueryRowContext(c, queryCount).Scan(&total)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	rows, err := s.db.QueryContext(c, queryData)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	defer rows.Close()
	var link model.Link

	for rows.Next() {
		err = rows.Scan(&link.ID, &link.UserID, &link.CreatedAt, &link.UpdatedAt, &link.Domain, &link.Keyword, &link.URL, &link.Title, &link.Active)

		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return
		}

		links = append(links, link)
	}

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return total, pages, links, e.ErrInternalServerError
	}

	pages = int(math.Ceil(float64(total) / float64(r.Limit)))

	return
}

// Update change specific link by ID.
func (s service) Update(c *gin.Context, payload model.Link) (err error) {
	log := logger.Logger(c)

	query := "UPDATE links SET title = $1 WHERE id = $2 AND user_id = $3"
	_, err = s.db.ExecContext(c, query, payload.Title, payload.ID, payload.UserID)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return e.ErrInternalServerError
	}

	return
}

// Delete delete a link by ID.
func (s service) Delete(c *gin.Context, linkID, userID string) (err error) {
	log := logger.Logger(c)

	link := model.Link{}

	query := "SELECT id, domain, keyword FROM links WHERE id = $1 AND user_id = $2 AND active = true"
	err = s.db.QueryRowContext(c, query, linkID, userID).Scan(&link.ID, &link.Domain, &link.Keyword)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Error().Caller().Msg(err.Error())
			return e.ErrInternalServerError
		}
	}

	if link.ID == "" {
		message := fmt.Sprintf("Link enable with ID = '%s' not found", linkID)
		log.Info().Caller().Msg(message)
		return e.ErrLinkNotFound
	}

	key := fmt.Sprintf("link_%s_%s", link.Domain, link.Keyword)
	_, err = s.cache.Del(c, key).Result()

	// Keep going on error from cache.
	if err != nil {
		log.Error().Caller().Msg(err.Error())
	}

	query = "UPDATE links SET active = false, updated_at = $1 WHERE id = $2 AND user_id = $3"
	_, err = s.db.ExecContext(c, query, time.Now(), linkID, userID)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return e.ErrInternalServerError
	}

	return
}
