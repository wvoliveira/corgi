package token

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/jwt"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	"gorm.io/gorm"
)

// Service encapsulates the authentication logic.
type Service interface {
	Refresh(ctx context.Context, token model.Token) (model.Token, model.Token, error)

	NewHTTP(r *mux.Router)
	HTTPRefresh(w http.ResponseWriter, r *http.Request)
}

type service struct {
	db              *gorm.DB
	secret          string
	tokenExpiration int
	store           *sessions.CookieStore
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB, secret string, tokenExpiration int, store *sessions.CookieStore) Service {
	return service{db, secret, tokenExpiration, store}
}

// Refresh authenticates a user and generates a new access and refresh JWT token if needed.
// Otherwise, an error is returned.
func (s service) Refresh(ctx context.Context, payload model.Token) (tokenAccess, tokenRefresh model.Token, err error) {
	l := logger.Logger(ctx)

	if err = s.db.Model(&model.Token{}).Where("id = ?", payload.ID).First(&tokenRefresh).Error; err != nil {
		l.Warn().Caller().Msg(err.Error())
		return tokenAccess, tokenRefresh, e.ErrUnauthorized
	}

	// Check refresh token from database.
	claims, ok := jwt.ValidToken(s.secret, tokenRefresh.Token)
	if !ok {
		l.Warn().Caller().Msg(err.Error())
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
		if err = s.db.Model(&model.Token{}).Create(&tokenRefresh).Error; err != nil {
			l.Error().Caller().Msg(err.Error())
			return tokenAccess, tokenRefresh, errors.New("error to create refresh token")
		}
		if err = s.db.Model(&model.Token{}).Where("id = ?", payload.ID).Delete(&payload).Error; err != nil {
			l.Error().Caller().Msg(err.Error())
			return tokenAccess, tokenRefresh, errors.New("error to delete refresh token from database")
		}
	}
	return
}
