package redirect

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/elga-io/corgi/pkg/queue"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	// Main business logic.
	FindByKeyword(ctx context.Context, domain, keyword string) (link entity.Link, err error)

	// HTTP transport user request/response.
	HTTPNewTransport(r *gin.Engine)
	HTTPFindByKeyword(c *gin.Context)

	// NATS transport broker.
	NATSNewTransport()
	NATSFindByKeyword(ctx context.Context)

	// MQ consumers.
	MQNewTransport(queue.MessageClient, ConsumerConfig) Consumer

	// Middlewares.
	MiddlewareMetric(log.Logger) gin.HandlerFunc
}

type service struct {
	logger  log.Logger
	db      *gorm.DB
	broker  *nats.EncodedConn
	queuer  queue.Queuer
	store   cookie.Store
	enforce *casbin.Enforcer
}

// NewService creates a new public service.
func NewService(logger log.Logger, db *gorm.DB, broker *nats.EncodedConn, queuer queue.Queuer, store cookie.Store, enforce *casbin.Enforcer) Service {
	return service{logger, db, broker, queuer, store, enforce}
}

// FindByKeyword get a shortener link from keyword.
func (s service) FindByKeyword(ctx context.Context, domain, keyword string) (l entity.Link, err error) {
	logger := s.logger.With(ctx, "domain", domain, "keyword", keyword)

	err = s.db.Model(&entity.Link{}).Where("domain = ? AND keyword = ?", domain, keyword).Take(&l).Error
	if err == gorm.ErrRecordNotFound {
		logger.Infof("the link domain '%s' and keyword '%s' not found", domain, keyword)
		return l, e.ErrLinkNotFound
	}

	if err != nil {
		logger.Errorf("an errors occurred to get keyword from database: %s", err.Error())
		return
	}
	return
}

// Log store a log metadata to database.
func (s service) Log(ctx context.Context, payload entity.LinkLog) (err error) {
	logger := s.logger.With(ctx)

	fmt.Println("Log service")
	fmt.Println(payload)

	err = s.db.Debug().Model(&entity.LinkLog{}).Create(&payload).Error
	if err != nil {
		logger.Errorf("an errors occurred to create link log: %s", err.Error())
		return
	}
	return
}
