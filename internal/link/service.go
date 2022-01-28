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
	logger := s.logger.With(ctx, "user_id", link.UserID)
	logger.Infof("add link with url short with domain '%s' and keyword '%s'", link.Domain, link.Keyword)

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
		logger.Warnf("domain '%s' with keyword '%s' already exists", link.Domain, link.Keyword)
		return l, e.ErrLinkAlreadyExists
	}
	logger.Errorf("error when creating a new shortener link, look: %s", err.Error())
	return l, e.ErrInternalServerError
}

// FindByID get a shortener link from ID.
func (s service) FindByID(ctx context.Context, linkID, userID string) (l entity.Link, err error) {
	logger := s.logger.With(ctx, "user_id", userID)
	logger.Infof("find link with id '%s'", linkID)

	err = s.db.Model(&entity.Link{}).Where("id = ? AND user_id = ?", linkID, userID).Take(&l).Error
	if err == gorm.ErrRecordNotFound {
		logger.Infof("the link id '%s' not found from user_id '%s'", linkID, userID)
		return l, e.ErrLinkNotFound
	} else if err == nil {
		return
	}
	logger.Errorf("oh crap, an errors occurred: %s", err.Error())
	return
}

// FindAll get a list of links from database.
func (s service) FindAll(ctx context.Context, offset, limit int, sort, userID string) (total int64, pages int, links []entity.Link, err error) {
	logger := s.logger.With(ctx, "user_id", userID)

	s.db.Model(&entity.Link{}).Where("user_id = ?", userID).Count(&total)
	err = s.db.Model(&entity.Link{}).Where("user_id = ?", userID).Offset(offset).Limit(limit).Order(sort).Find(&links).Error

	if err == gorm.ErrRecordNotFound {
		logger.Infof("the links with '%d' offset and '%d' limit not found from user_id '%s'", offset, limit, userID)
		return total, pages, links, e.ErrLinkNotFound
	} else if err != nil {
		logger.Errorf("oh crap, an errors occurred: %s", err.Error())
		return
	}

	pages = int(math.Ceil(float64(total) / float64(limit)))
	return
}

// Update change specific link by ID.
func (s service) Update(ctx context.Context, link entity.Link) (l entity.Link, err error) {
	logger := s.logger.With(ctx, "user_id", link.UserID)
	logger.Infof("updating link with id '%s'", link.ID)

	link.UpdatedAt = time.Now()

	err = s.db.Model(&entity.Link{}).Where("id = ? AND user_id = ?", link.ID, link.UserID).Updates(&link).Error
	if err == gorm.ErrRecordNotFound {
		logger.Infof("the link id '%s' not found from user_id '%s'", link.ID, link.UserID)
		return l, e.ErrLinkNotFound
	} else if err == nil {
		l.ID = link.ID
		l.CreatedAt = link.CreatedAt
		l.UpdatedAt = link.UpdatedAt
		l.Domain = link.Domain
		l.Keyword = link.Keyword
		l.URL = link.URL
		l.Title = link.Title
		l.Active = link.Active
		return
	}

	logger.Errorf("oh crap, an errors occurred: %s", err.Error())
	return
}

// Delete delete a link by ID.
func (s service) Delete(ctx context.Context, linkID, userID string) (err error) {
	logger := s.logger.With(ctx, "user_id", userID)
	logger.Infof("deleting link with id '%s'", linkID)

	err = s.db.Debug().
		Model(&entity.Link{}).
		Clauses(clause.Returning{}).
		Where("id = ? AND user_id = ?", linkID, userID).
		Delete(&entity.Link{ID: linkID, UserID: userID}).Error

	if err == gorm.ErrRecordNotFound {
		logger.Infof("the link id '%s' not found from user_id '%s'", linkID, userID)
		return e.ErrLinkNotFound
	} else if err == nil {
		return
	}
	logger.Errorf("oh crap, an errors occurred: %s", err.Error())
	return
}
