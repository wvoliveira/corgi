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
	FindByID(*gin.Context, string, string) (model.User, error)
	UpdateByID(*gin.Context, string, string, string) error
	FindByUsername(*gin.Context, string, string) (model.User, error)

	NewHTTP(*gin.RouterGroup)
	HTTPFindByID(*gin.Context)
	HTTPUpdateByID(*gin.Context)
	HTTPFindByUsername(*gin.Context)
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
func (s service) FindByID(c *gin.Context, whoID string, id string) (user model.User, err error) {
	log := logger.Logger(c)

	query := "SELECT id, created_at, updated_at, username, name, role, active FROM users WHERE id = $1"

	err = s.db.QueryRowContext(c, query, id).Scan(
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
func (s service) UpdateByID(c *gin.Context, whoID, id, name string) (err error) {
	log := logger.Logger(c.Request.Context())

	// TODO:
	// 	- check if user is updating yourself
	// 	- update all values
	query := "UPDATE users SET name = $1 WHERE id = $2"
	_, err = s.db.ExecContext(c, query, name, id)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	return
}

// Find get a user from username.
func (s service) FindByUsername(c *gin.Context, whoID string, username string) (user model.User, err error) {
	log := logger.Logger(c)

	query := "SELECT id, created_at, updated_at, username, name, role, active FROM users WHERE username = $1"

	err = s.db.QueryRowContext(c, query, username).Scan(
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
