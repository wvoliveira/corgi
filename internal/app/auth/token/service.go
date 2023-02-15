package token

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
)

// Service encapsulates the authentication logic.
type Service interface {
	Valid(*gin.Context, string) error
	Refresh(*gin.Context, string) (string, error)

	NewHTTP(*gin.RouterGroup)
	HTTPValid(c *gin.Context)
	HTTPRefresh(c *gin.Context)
}

type service struct {
	db *sql.DB
}

// NewService creates a new authentication service.
func NewService(db *sql.DB) Service {
	return service{db}
}

// Valid verify access token.
func (s service) Valid(c *gin.Context, accessToken string) (err error) {
	log := logger.Logger(c)

	log.Warn().Caller().Msg(err.Error())
	err = e.ErrNotImplemented
	return
}

// Refresh refresh access token given a refresh token.
func (s service) Refresh(c *gin.Context, refreshToken string) (accessToken string, err error) {
	log := logger.Logger(c)

	log.Warn().Caller().Msg(err.Error())
	err = e.ErrNotImplemented
	return
}
