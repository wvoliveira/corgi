package health

import (
	"context"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gorilla/mux"
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
	store   cookie.Store
	enforce *casbin.Enforcer
	version string
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB, secret string, store cookie.Store, enforce *casbin.Enforcer, version string) Service {
	return service{db, secret, store, enforce, version}
}

// Health create a new shortener link.
func (s service) Health(_ context.Context) (err error) {
	// l := log.Ctx(ctx)
	return
}
