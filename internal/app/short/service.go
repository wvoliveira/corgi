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

	key := fmt.Sprintf("link_%s_%s", domain, keyword)

	log.Info().Caller().Msg("Trying to collect data from cache")
	val, err := s.cache.Get(c, key).Result()

	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Error().Caller().Msg(err.Error())
		}
	}

	if val != "" {
		log.Info().Caller().Msg("OK, I gotten from cache!")
		m.URL = val
		return
	}

	log.Info().Caller().Msg("No, link is not cached!")

	query := "SELECT url FROM links WHERE domain = $1 AND keyword = $2 AND active = true"
	err = s.db.QueryRowContext(c, query, domain, keyword).Scan(&m.URL)

	if errors.Is(err, sql.ErrNoRows) {
		return m, e.ErrLinkNotFound
	}

	status := s.cache.Set(c, key, m.URL, 10*time.Minute)
	err = status.Err()

	if err != nil {
		fmt.Println(err.Error())
	}

	return
}
