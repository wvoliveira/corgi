package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"gorm.io/gorm"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Find(ctx context.Context, userID string) (entity.User, error)
	Update(ctx context.Context, user entity.User) (entity.User, error)

	NewHTTP(r *mux.Router)
	HTTPFind(w http.ResponseWriter, r *http.Request)
	HTTPUpdate(w http.ResponseWriter, r *http.Request)
}

type service struct {
	db     *gorm.DB
	secret string
	store  *sessions.CookieStore
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB, secret string, store *sessions.CookieStore) Service {
	return service{db, secret, store}
}

// Find get a shortener link from ID.
func (s service) Find(ctx context.Context, userID string) (user entity.User, err error) {
	l := logger.Logger(ctx)

	if userID == "anonymous" {
		user.Name = "Anonymous"
		return
	}

	user.ID = userID
	err = s.db.
		Model(&user).
		Preload("Identities").
		Find(&user).Error

	if err == gorm.ErrRecordNotFound {
		l.Info().Caller().Msg(fmt.Sprintf("the user with user_id '%s' not found", userID))
		return user, e.ErrUserNotFound
	} else if err == nil {
		return
	}

	l.Error().Caller().Msg(err.Error())
	return
}

// Update change specific link by ID.
func (s service) Update(ctx context.Context, reqUser entity.User) (user entity.User, err error) {
	l := logger.Logger(ctx)

	if reqUser.ID == "anonymous" {
		return user, e.ErrUnauthorized
	}

	reqUser.UpdatedAt = time.Now()
	err = s.db.Model(&entity.User{}).
		Where("id = ?", reqUser.ID).
		Updates(&reqUser).Error

	if err == gorm.ErrRecordNotFound {
		l.Info().Caller().Msg(fmt.Sprintf("the user with id '%s' not found", reqUser.ID))
		return user, e.ErrUserNotFound
	} else if err == nil {
		user = reqUser
		return
	}

	l.Error().Caller().Msg(err.Error())
	return
}
