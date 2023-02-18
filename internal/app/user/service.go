package user

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Find(*gin.Context, string) (model.User, error)
	Update(*gin.Context, model.User) error

	NewHTTP(*gin.RouterGroup)
	HTTPFind(*gin.Context)
	HTTPUpdate(*gin.Context)
}

type service struct {
	db    *sql.DB
	cache *redis.Client
}

// NewService creates a new user management service.
func NewService(db *sql.DB, cache *redis.Client) Service {
	return service{db, cache}
}

// Find get a shortener link from ID.
func (s service) Find(c *gin.Context, userID string) (user model.User, err error) {
	log := logger.Logger(c)

	if userID == "0" {
		user.Name = "Anonymous"
		return
	}

	query := "SELECT id, created_at, updated_at, username, name, role, active FROM users WHERE id = $1"
	err = s.db.QueryRowContext(c, query, userID).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Username, &user.Name, &user.Role, &user.Active)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return user, e.ErrUserNotFound
	}

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	return
}

// Update change specific link by ID.
func (s service) Update(c *gin.Context, req model.User) (err error) {
	log := logger.Logger(c.Request.Context())

	if req.ID == "0" {
		return e.ErrUnauthorized
	}

	// TODO: update all values
	query := "UPDATE users SET name = $1 WHERE id = $2"
	_, err = s.db.ExecContext(c, query, req.Name, req.ID)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	return
}
