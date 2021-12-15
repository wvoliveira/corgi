package public

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	location2 "github.com/elga-io/corgi/pkg/location"
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
			go func() {
				ipLong := location2.IPv4ToLong(linkLog.RemoteAddress)
				location := entity.LocationIPv4{}
				if err = s.db.Debug().Model(&entity.LocationIPv4{}).Where("range_start <= ? AND range_end >= ?", ipLong, ipLong).Take(&location).Error; err != nil {
					logger.Warnf("error to get ipv4 location: %s", err.Error())
					return
				}

				linkLog.LinkID = l.ID
				linkLog.LocationIPv4ID = location.ID
				err = s.db.Debug().Model(entity.LinkLog{}).Create(&linkLog).Error
				if err != nil {
					logger.Warnf("error to create link log: %s", err.Error())
				}
			}()
		}
		return
	}
	logger.Errorf("oh crap, an errors occurred: %s", err.Error())
	return
}
