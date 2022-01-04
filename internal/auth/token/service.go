package token

import (
	"context"
	"errors"
	"github.com/casbin/casbin/v2"
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/jwt"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

// Service encapsulates the authentication logic.
type Service interface {
	Refresh(ctx context.Context, token entity.Token) (entity.Token, entity.Token, error)
	HTTPRefresh(c *gin.Context)
	Routers(r *gin.Engine)
}

type service struct {
	logger          log.Logger
	db              *gorm.DB
	secret          string
	tokenExpiration int
	store           cookie.Store
	enforce         *casbin.Enforcer
}

// NewService creates a new authentication service.
func NewService(logger log.Logger, db *gorm.DB, secret string, tokenExpiration int, store cookie.Store, enforce *casbin.Enforcer) Service {
	return service{logger, db, secret, tokenExpiration, store, enforce}
}

// Refresh authenticates a user and generates a new access and refresh JWT token if needed.
// Otherwise, an error is returned.
func (s service) Refresh(ctx context.Context, payload entity.Token) (tokenAccess, tokenRefresh entity.Token, err error) {
	logger := s.logger.With(ctx)

	if err = s.db.Debug().Model(&entity.Token{}).Where("id = ?", payload.ID).First(&tokenRefresh).Error; err != nil {
		logger.Warnf("error to get refresh token from our database: %s", err.Error())
		return tokenAccess, tokenRefresh, e.ErrUnauthorized
	}

	// Check refresh token from database.
	claims, ok := jwt.ValidToken(s.secret, tokenRefresh.Token)
	if !ok {
		logger.Warnf("invalid refresh token from our database: %s", err.Error())
		// TODO: Delete from database if not valid.
		return tokenAccess, tokenRefresh, e.ErrUnauthorized
	}

	exp := claims["exp"].(int64)
	tm := time.Unix(exp, 0)
	remains := time.Since(tm).Hours()

	genRefresh := false
	// If there is 2 hours left, create a new refresh token.
	if remains < -2 {
		genRefresh = true
	}

	tokenAccess, err = jwt.UpdateAccessToken(s.secret, claims)
	if err != nil {
		return tokenAccess, tokenRefresh, errors.New("error to generate access token: " + err.Error())
	}

	tokenRefresh, err = jwt.UpdateRefreshToken(s.secret, claims)
	if err != nil {
		return tokenAccess, tokenRefresh, errors.New("error to generate refresh token: " + err.Error())
	}

	if genRefresh {
		if err = s.db.Debug().Model(&entity.Token{}).Create(&tokenRefresh).Error; err != nil {
			logger.Error("error to create refresh token", err.Error())
			return tokenAccess, tokenRefresh, errors.New("error to create refresh token")
		}
		if err = s.db.Debug().Model(&entity.Token{}).Where("id = ?", payload.ID).Delete(&payload).Error; err != nil {
			logger.Error("error to delete refresh token from database", err.Error())
			return tokenAccess, tokenRefresh, errors.New("error to delete refresh token from database")
		}
	}
	return
}
