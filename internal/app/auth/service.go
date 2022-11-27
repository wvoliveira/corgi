package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"gorm.io/gorm"
)

// Service encapsulates the authentication logic.
type Service interface {
	Logout(c *gin.Context, token entity.Token) error

	NewHTTP(*gin.RouterGroup)
	HTTPLogout() gin.HandlerFunc
}

type service struct {
	db *gorm.DB
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB) Service {
	return service{db}
}

// Logout remove cookie and refresh token from database.
func (s service) Logout(c *gin.Context, token entity.Token) (err error) {
	l := logger.Logger(c.Request.Context())

	if token.UserID == "anonymous" {
		return e.ErrUnauthorized
	}

	err = s.db.Debug().Model(&entity.Token{}).Where("id = ? AND user_id = ?", token.ID, token.UserID).Delete(&token).Error
	if err != nil {
		l.Error().Caller().Msg(err.Error())
		return
	}
	return
}
