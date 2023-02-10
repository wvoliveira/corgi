package user

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/gin-gonic/gin"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Find(*gin.Context, string) (model.User, error)
	Update(*gin.Context, model.User) (model.User, error)

	NewHTTP(*gin.RouterGroup)
	HTTPFind(*gin.Context)
	HTTPUpdate(*gin.Context)
}

type service struct {
	db *sql.DB
	kv *badger.DB
}

// NewService creates a new user management service.
func NewService(db *sql.DB, kv *badger.DB) Service {
	return service{db, kv}
}

// Find get a shortener link from ID.
func (s service) Find(c *gin.Context, userID string) (user model.User, err error) {
	log := logger.Logger(c.Request.Context())

	if userID == "anonymous" {
		user.Name = "Anonymous"
		return
	}

	query := "SELECT * FROM users WHERE id = ?"
	err = s.db.QueryRowContext(c, query, userID).Scan(user)

	if err != nil {
		log.Info().Caller().Msg(fmt.Sprintf("the user with user_id \"%s\" was not found", userID))
		return user, e.ErrUserNotFound
	}

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	return
}

// Update change specific link by ID.
func (s service) Update(c *gin.Context, req model.User) (user model.User, err error) {
	log := logger.Logger(c.Request.Context())

	if req.ID == "anonymous" {
		return user, e.ErrUnauthorized
	}

	// TODO: update all values
	query := "UPDATE users SET updated_at = ?, name = ? WHERE id = ?"
	_, err = s.db.ExecContext(c, query, time.Now(), req.Name, req.ID)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	return
}
