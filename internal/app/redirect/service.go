package redirect

import (
	"database/sql"

	"github.com/dgraph-io/badger"
	"github.com/gin-gonic/gin"
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
	db *sql.DB
	kv *badger.DB
}

// NewService creates a new public service.
func NewService(db *sql.DB, kv *badger.DB) Service {
	return service{db, kv}
}

// Find get a shortener link from keyword.
func (s service) Find(c *gin.Context, domain, keyword string) (m model.Link, err error) {
	log := logger.Logger(c)

	query := "SELECT * FROM links WHERE domain = ? AND keyword = ?"
	err = s.db.QueryRowContext(c, query, domain, keyword).Scan(&m)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return m, e.ErrInternalServerError
	}

	// link := fmt.Sprintf("%s/%s", domain, keyword)
	// go increaseClick(c, s.kv, link, time.Now())

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
