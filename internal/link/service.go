package link

import (
	"context"
	"errors"
	"github.com/casbin/casbin/v2"
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
	"time"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Add(ctx context.Context, payload entity.Link) (link entity.Link, err error)
	FindByID(ctx context.Context, linkID, userID string) (link entity.Link, err error)
	FindAll(ctx context.Context, offset, limit int, sort, userID string) (total int64, pages int, links []entity.Link, err error)
	Update(ctx context.Context, payload entity.Link) (link entity.Link, err error)
	Delete(ctx context.Context, linkID, userID string) (err error)

	HTTPNewTransport(r *gin.Engine)
	HTTPAdd(c *gin.Context)
	HTTPFindByID(c *gin.Context)
	HTTPFindAll(c *gin.Context)
	HTTPUpdate(c *gin.Context)
	HTTPDelete(c *gin.Context)

	NatsNewTransport()
	NatsAdd(ctx context.Context)
	NatsFindByID(ctx context.Context)
	NatsFindAll(ctx context.Context)
	NatsUpdate(ctx context.Context)
	NatsDelete(ctx context.Context)
}

type service struct {
	logger  log.Logger
	db      *gorm.DB
	broker  *nats.EncodedConn
	secret  string
	store   cookie.Store
	enforce *casbin.Enforcer
}

// NewService creates a new authentication service.
func NewService(logger log.Logger, db *gorm.DB, broker *nats.EncodedConn, secret string, store cookie.Store, enforce *casbin.Enforcer) Service {
	return service{logger, db, broker, secret, store, enforce}
}

// Add create a new shortener link.
func (s service) Add(ctx context.Context, link entity.Link) (l entity.Link, err error) {
	logger := s.logger.With(ctx, "user_id", link.UserID, "link_domain", link.Domain, "link_keyword", link.Keyword)

	if err = checkLink(logger, link); err != nil {
		return
	}

	err = s.db.Model(&entity.Link{}).Where("domain = ? AND keyword = ?", link.Domain, link.Keyword).Take(&l).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		l.ID = uuid.New().String()
		l.CreatedAt = time.Now()
		l.Domain = link.Domain
		l.Keyword = link.Keyword
		l.URL = link.URL
		l.Title = link.Title
		l.Active = "true"
		l.UserID = link.UserID

		err = s.db.Model(&entity.Link{}).Create(&l).Error
		return l, err
	} else if err == nil {
		logger.Warnf("domain with keyword already exists")
		return l, e.ErrLinkAlreadyExists
	}
	logger.Error(err.Error())
	return l, e.ErrInternalServerError
}

// FindByID get a shortener link from ID.
func (s service) FindByID(ctx context.Context, linkID, userID string) (l entity.Link, err error) {
	logger := s.logger.With(ctx, "user_id", userID, "link_id", linkID)

	err = s.db.Model(&entity.Link{}).Where("id = ? AND user_id = ?", linkID, userID).Take(&l).Error
	if err == gorm.ErrRecordNotFound {
		logger.Info("link not found")
		return l, e.ErrLinkNotFound
	} else if err == nil {
		return
	}
	logger.Error(err.Error())
	return
}

// FindAll get a list of links from database.
func (s service) FindAll(ctx context.Context, offset, limit int, sort, userID string) (total int64, pages int, links []entity.Link, err error) {
	logger := s.logger.With(ctx, "user_id", userID, "offset", offset, "limit", limit, "sort", sort)

	err = s.db.Model(&entity.Link{}).Where("user_id = ?", userID).Count(&total).Offset(offset).Limit(limit).Order(sort).Find(&links).Error

	if err == gorm.ErrRecordNotFound {
		logger.Info("links not found", offset, limit, userID)
		return total, pages, links, e.ErrLinkNotFound
	} else if err != nil {
		logger.Error(err.Error())
		return
	}

	pages = int(math.Ceil(float64(total) / float64(limit)))
	return
}

// Update change specific link by ID.
func (s service) Update(ctx context.Context, link entity.Link) (l entity.Link, err error) {
	logger := s.logger.With(ctx, "id", link.ID, "user_id", link.UserID)

	err = s.db.Model(&entity.Link{}).Where("id = ? AND user_id = ?", link.ID, link.UserID).Updates(&link).First(&l).Error
	if err == gorm.ErrRecordNotFound {
		logger.Info("link not found")
		return l, e.ErrLinkNotFound
	} else if err != nil {
		return l, e.ErrInternalServerError
	}
	return
}

// Delete delete a link by ID.
func (s service) Delete(ctx context.Context, linkID, userID string) (err error) {
	logger := s.logger.With(ctx, "user_id", userID, "link_id", linkID)

	err = s.db.Debug().
		Model(&entity.Link{}).
		Clauses(clause.Returning{}).
		Where("id = ? AND user_id = ?", linkID, userID).
		Delete(&entity.Link{ID: linkID, UserID: userID}).Error

	if err == gorm.ErrRecordNotFound {
		logger.Info("link not found")
		return e.ErrLinkNotFound
	} else if err == nil {
		return
	}
	logger.Error(err.Error())
	return
}
