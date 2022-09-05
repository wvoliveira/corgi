package auth

import (
	"context"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/wvoliveira/corgi/internal/app/entity"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"gorm.io/gorm"
)

// Service encapsulates the authentication logic.
type Service interface {
	Logout(ctx context.Context, token entity.Token) error

	NewHTTP(r *mux.Router)
	HTTPLogout(w http.ResponseWriter, r *http.Request)
}

type service struct {
	db      *gorm.DB
	secret  string
	store   *sessions.CookieStore
	enforce *casbin.Enforcer
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB, secret string, store *sessions.CookieStore, enforce *casbin.Enforcer) Service {
	return service{db, secret, store, enforce}
}

// Logout remove cookie and refresh token from database.
func (s service) Logout(ctx context.Context, token entity.Token) (err error) {
	l := logger.Logger(ctx)

	err = s.db.Debug().Model(&entity.Token{}).Where("id = ? AND user_id = ?", token.ID, token.UserID).Delete(&token).Error
	if err != nil {
		l.Error().Caller().Msg(err.Error())
		return
	}
	return
}
