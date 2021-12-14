package public

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	FindByKeyword(ctx context.Context, domain, keyword string, linkLog entity.LinkLog, unique bool) (link entity.Link, err error)
	HTTPFindByKeyword(c *gin.Context)
	Routers(r *gin.Engine)
}

type service struct {
	logger  log.Logger
	db      *gorm.DB
	store   cookie.Store
	enforce *casbin.Enforcer
}

// NewService creates a new public service.
func NewService(logger log.Logger, db *gorm.DB, store cookie.Store, enforce *casbin.Enforcer) Service {
	return service{logger, db, store, enforce}
}

// FindByKeyword get a shortener link from keyword.
func (s service) FindByKeyword(ctx context.Context, domain, keyword string, linkLog entity.LinkLog, unique bool) (l entity.Link, err error) {
	logger := s.logger.With(ctx, "domain", domain, "keyword", keyword)

	err = checkLink(logger, entity.Link{Domain: domain, Keyword: keyword})
	if err != nil {
		return
	}

	err = s.db.Model(&entity.Link{}).Where("domain = ? AND keyword = ?", domain, keyword).Take(&l).Error
	if err == gorm.ErrRecordNotFound {
		logger.Infof("the link domain '%s' and keyword '%s' not found", domain, keyword)
		return l, e.ErrLinkNotFound
	} else if err == nil {
		if unique {
			linkLog.LinkID = l.ID
			go s.db.Debug().Model(entity.LinkLog{}).Create(&linkLog)
		}
		return
	}
	logger.Errorf("oh crap, an errors occurred: %s", err.Error())
	return
}
