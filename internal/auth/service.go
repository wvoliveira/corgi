package auth

import (
	"context"
	"github.com/elga-io/corgi/internal/entity"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Service encapsulates the authentication logic.
type Service interface {
	Logout(ctx context.Context, token entity.Token) error
	HTTPLogout(c *gin.Context)
	Routers(r *gin.Engine)
}

type service struct {
	logger log.Logger
	db     *gorm.DB
	secret string
	store  cookie.Store
}

// NewService creates a new authentication service.
func NewService(logger log.Logger, db *gorm.DB, secret string, store cookie.Store) Service {
	return service{logger, db, secret, store}
}

// Logout remove cookie and refresh token from database.
func (s service) Logout(ctx context.Context, token entity.Token) (err error) {
	logger := s.logger.With(ctx, "user_id", token.UserID)
	logger.Info("Logout user deleting cookie and refresh token.")

	err = s.db.Debug().Model(&entity.Token{}).Delete(&token).Error
	if err != nil {
		logger.Errorf("error to delete token in database: %s", err.Error())
		return
	}
	return
}
