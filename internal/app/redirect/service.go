package redirect

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"gorm.io/gorm"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Find(ctx context.Context, domain, keyword string) (link entity.Link, err error)

	NewHTTP(r *mux.Router)
	HTTPFind(w http.ResponseWriter, r *http.Request)
}

type service struct {
	db    *gorm.DB
	store *sessions.CookieStore
}

// NewService creates a new public service.
func NewService(db *gorm.DB, store *sessions.CookieStore) Service {
	return service{db, store}
}

// Find get a shortener link from keyword.
func (s service) Find(ctx context.Context, domain, keyword string) (li entity.Link, err error) {
	l := logger.Logger(ctx)

	err = s.db.Model(&entity.Link{}).Where("domain = ? AND keyword = ?", domain, keyword).Take(&li).Error
	if err == gorm.ErrRecordNotFound {
		l.Info().Caller().Msg(fmt.Sprintf("the link domain '%s' and keyword '%s' not found", domain, keyword))
		return li, e.ErrLinkNotFound
	}

	if err != nil {
		l.Error().Caller().Msg(err.Error())
		return
	}
	return
}

// Log store a log metadata to database.
func (s service) Log(ctx context.Context, payload entity.LinkLog) (err error) {
	l := logger.Logger(ctx)

	err = s.db.Debug().Model(&entity.LinkLog{}).Create(&payload).Error
	if err != nil {
		l.Error().Caller().Msg(err.Error())
		return
	}
	return
}
