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
	Refresh(ctx context.Context, token entity.Token) (entity.Token, error)
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
func (s service) Refresh(ctx context.Context, payload entity.Token) (token entity.Token, err error) {
	logger := s.logger.With(ctx)

	claims, ok := jwt.ValidToken(s.secret, payload.RefreshToken)
	if !ok {
		logger.Error("invalid refresh token", err.Error())
		return token, e.ErrUnauthorized
	}

	token.ID = claims["id"].(string)
	if err = s.db.Debug().Model(&entity.Token{}).Where("id = ?", payload.ID).First(&token).Error; err != nil {
		logger.Error("error to get refresh token from our database", err.Error())
		return token, e.ErrUnauthorized
	}

	if claims, ok = jwt.ValidToken(s.secret, token.RefreshToken); !ok {
		logger.Errorf("invalid refresh token from our database", err.Error())
		// TODO: Delete from database if not valid.
		return token, e.ErrUnauthorized
	}

	exp := claims["exp"].(int64)
	tm := time.Unix(exp, 0)
	remains := time.Since(tm).Hours()

	genRefresh := false
	// If there is 2 hours left, create a new refresh token.
	if remains < 2 {
		genRefresh = true
	}

	accessToken, refreshToken, err := jwt.GenerateTokens(s.secret, claims, genRefresh)
	if err != nil {
		return token, errors.New("error to generate access token: " + err.Error())
	}

	if genRefresh {
		if err = s.db.Debug().Model(&entity.Token{}).Create(&refreshToken).Error; err != nil {
			logger.Error("error to create refresh token", err.Error())
			return token, errors.New("error to create refresh token")
		}
		if err = s.db.Debug().Model(&entity.Token{}).Where("id = ?", payload.ID).Delete(&token).Error; err != nil {
			logger.Error("error to delete refresh token from database", err.Error())
			return token, errors.New("error to delete refresh token from database")
		}
		token.RefreshToken = refreshToken.RefreshToken
	}

	token.AccessToken = accessToken.AccessToken
	token.AccessExpires = accessToken.AccessExpires
	return
}
