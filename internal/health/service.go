package health

import (
	"context"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	health(ctx context.Context) error
	//Healthcheck(ctx context.Context) error

	httpHealth(c *gin.Context)
	//HTTPHealthcheck(c *gin.Context)

	Routers(r *gin.RouterGroup)
}

type service struct {
	logger  log.Logger
	db      *gorm.DB
	version string
}

// NewService creates a new authentication service.
func NewService(logger log.Logger, db *gorm.DB, version string) Service {
	return service{logger, db, version}
}

// Health create a new shortener link.
func (s service) health(_ context.Context) (err error) {
	//logger := s.logger.With(ctx)
	return
}
