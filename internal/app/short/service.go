package short

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Find(*gin.Context, string, string) (model.Link, error)

	NewHTTP(*gin.RouterGroup)
	HTTPFind(*gin.Context)
}

type service struct {
	db    *sql.DB
	cache *redis.Client
}

// NewService creates a new public service.
func NewService(db *sql.DB, cache *redis.Client) Service {
	return service{db, cache}
}

// Find get a shortener link from keyword.
func (s service) Find(c *gin.Context, domain, keyword string) (m model.Link, err error) {
	log := logger.Logger(c)

	log.Info().Caller().Msg("Trying to collect data from cache")
	key := fmt.Sprintf("link_%s_%s", domain, keyword)
	val, err := s.cache.Get(c, key).Result()

	if err != nil {
		log.Error().Caller().Msg(err.Error())
	} else {
		log.Info().Caller().Msg("OK, I gotten from cache!")
		m.URL = val
		return
	}

	query := "SELECT url FROM links WHERE domain = $1 AND keyword = $2"
	err = s.db.QueryRowContext(c, query, domain, keyword).Scan(&m.URL)

	if errors.Is(err, sql.ErrNoRows) {
		log.Error().Caller().Msg(err.Error())
		return m, e.ErrLinkNotFound
	}

	status := s.cache.Set(c, key, m.URL, 10*time.Minute)
	err = status.Err()

	if err != nil {
		fmt.Println(err.Error())
	}

	return
}

// Log store a log metadata to database.
func (s service) Log(c *gin.Context, payload model.LinkLog) (err error) {
	// log := logger.Logger(c)

	// err = s.db.Debug().Model(&model.LinkLog{}).Create(&payload).Error

	// if err != nil {
	// 	log.Error().Caller().Msg(err.Error())
	// 	return e.ErrInternalServerError
	// }

	return
}
