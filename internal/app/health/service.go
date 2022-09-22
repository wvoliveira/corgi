package health

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Health(ctx context.Context) error

	NewHTTP(r *mux.Router)
	HTTPHealth(w http.ResponseWriter, r *http.Request)
}

type service struct {
	db      *gorm.DB
	secret  string
	store   *sessions.CookieStore
	version string
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB, secret string, store *sessions.CookieStore, version string) Service {
	return service{db, secret, store, version}
}

// Health create a new shortener link.
func (s service) Health(_ context.Context) (err error) {
	// l := logger.Logger(ctx)
	return
}
