package health

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	health(ctx context.Context) error
	//Healthcheck(ctx context.Context) error

	httpHealth(c *gin.Context)
	//HTTPHealthcheck(c *gin.Context)

	Routers(r *gin.Engine)
}

type service struct {
	logger  log.Logger
	db      *gorm.DB
	secret  string
	store   cookie.Store
	enforce *casbin.Enforcer
	version string
}

// NewService creates a new authentication service.
func NewService(logger log.Logger, db *gorm.DB, secret string, store cookie.Store, enforce *casbin.Enforcer, version string) Service {
	return service{logger, db, secret, store, enforce, version}
}

// Health create a new shortener link.
func (s service) health(_ context.Context) (err error) {
	//logger := s.logger.With(ctx)
	return
}
