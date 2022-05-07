package health

import (
	"context"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Health(ctx context.Context) error

	NewHTTP(r *gin.Engine)
	HTTPHealth(c *gin.Context)
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
