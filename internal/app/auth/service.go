package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"gorm.io/gorm"
)

// Service encapsulates the authentication logic.
type Service interface {
	Logout(*gin.Context, entity.User) error

	NewHTTP(*gin.RouterGroup)
	HTTPLogout(c *gin.Context)
}

type service struct {
	db *gorm.DB
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB) Service {
	return service{db}
}

// Logout remove cookie and refresh token from database.
func (s service) Logout(c *gin.Context, user entity.User) (err error) {
	log := logger.Logger(c.Request.Context())

	// TODO: make something with user logout session.
	log.Info().Caller().Msg("user logout but do nothing in backend yet")

	return
}
