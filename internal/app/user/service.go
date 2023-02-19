package user

import (
	"database/sql"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	FindMe(*gin.Context, string) (model.User, error)
	UpdateMe(*gin.Context, string, string) error
	FindByIDorUsername(*gin.Context, string, string) (model.User, error)
	UpdateByIDorUsername(*gin.Context, string, string, string) error

	NewHTTP(*gin.RouterGroup)
	HTTPFindByIDorUsername(*gin.Context)
	HTTPUpdateByIDorUsername(*gin.Context)
}

type service struct {
	db    *sql.DB
	cache *redis.Client
}

// NewService creates a new user management service.
func NewService(db *sql.DB, cache *redis.Client) Service {
	return service{db, cache}
}

// FindMe get my personal info.
func (s service) FindMe(c *gin.Context, whoID string) (user model.User, err error) {
	log := logger.Logger(c)

	query := "SELECT id, created_at, updated_at, username, name, role, active FROM users WHERE id = $1"
	err = s.db.QueryRowContext(c, query, whoID).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Username, &user.Name, &user.Role, &user.Active)

	if err != nil {
		log.Error().Caller().Msg(err.Error())

		if errors.Is(err, sql.ErrNoRows) {
			return user, e.ErrUserNotFound
		}

		return
	}

	return
}

// Update change my profile.
func (s service) UpdateMe(c *gin.Context, whoID string, name string) (err error) {
	log := logger.Logger(c.Request.Context())

	// TODO: update all values
	query := "UPDATE users SET name = $1 WHERE id = $2"

	_, err = s.db.ExecContext(c, query, name, whoID)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	return
}

// Find get a shortener link from ID or username.
func (s service) FindByIDorUsername(c *gin.Context, whoID string, idOrUsername string) (user model.User, err error) {
	log := logger.Logger(c)

	query := "SELECT id, created_at, updated_at, username, name, role, active FROM users WHERE id = $1 OR username = $1"

	err = s.db.QueryRowContext(c, query, idOrUsername).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Username, &user.Name, &user.Role, &user.Active)

	if err != nil {
		log.Error().Caller().Msg(err.Error())

		if errors.Is(err, sql.ErrNoRows) {
			return user, e.ErrUserNotFound
		}

		return
	}

	return
}

// Update change specific link by ID or username.
func (s service) UpdateByIDorUsername(c *gin.Context, whoID, idOrUsername, name string) (err error) {
	log := logger.Logger(c.Request.Context())

	// TODO:
	// 	- check if user is updating yourself
	// 	- update all values
	query := "UPDATE users SET name = $1 WHERE id = $2 or username = $2"
	_, err = s.db.ExecContext(c, query, name, idOrUsername)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	return
}
