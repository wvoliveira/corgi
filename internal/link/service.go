package link

import (
	"context"
	"errors"
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Add(ctx context.Context, link addRequest) (entity.Link, error)
	FindByID(ctx context.Context, link findByIDRequest) (entity.Link, error)
	FindAll(ctx context.Context, link findAllRequest) ([]entity.Link, error)
	Update(ctx context.Context, link updateRequest) (entity.Link, error)
	Delete(ctx context.Context, link deleteRequest) error

	HTTPAdd(c *gin.Context)
	HTTPFindByID(c *gin.Context)
	HTTPFindAll(c *gin.Context)
	HTTPUpdate(c *gin.Context)
	HTTPDelete(c *gin.Context)

	Routers(r *gin.RouterGroup)
}

type service struct {
	logger log.Logger
	db     *gorm.DB
}

// NewService creates a new authentication service.
func NewService(logger log.Logger, db *gorm.DB) Service {
	return service{logger, db}
}

// Add create a new shortener link.
func (s service) Add(ctx context.Context, link addRequest) (l entity.Link, err error) {
	logger := s.logger.With(ctx, "user_id", link.UserID)
	logger.Infof("add link with url short '%s'", link.URLShort)

	err = s.db.Model(&entity.Link{}).Where("url_short = ?", link.URLShort).Take(&l).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		l.ID = uuid.New().String()
		l.CreatedAt = time.Now()
		l.URLShort = link.URLShort
		l.URLFull = link.URLFull
		l.Title = link.Title
		l.Active = "true"
		l.UserID = link.UserID

		err = s.db.Model(&entity.Link{}).Create(&l).Error
		return l, err
	} else if err == nil {
		logger.Warnf("shortener link '%s' already exists", link.URLShort)
		return l, e.ErrLinkKeywordAlreadyExists
	}
	logger.Errorf("error when creating a new shortener link, look: %s", err.Error())
	return l, e.ErrInternalServerError
}

// FindByID get a shortener link from ID.
func (s service) FindByID(ctx context.Context, link findByIDRequest) (l entity.Link, err error) {
	logger := s.logger.With(ctx, "user_id", link.UserID)
	logger.Infof("find link with id '%s'", link.ID)

	err = s.db.Model(&entity.Link{}).Where("id = ? AND user_id = ?", link.ID, link.UserID).First(&l).Error
	if err == gorm.ErrRecordNotFound {
		logger.Infof("the link id '%s' not found from user_id '%s'", link.ID, link.UserID)
		return l, e.ErrLinkNotFound
	} else if err == nil {
		return
	}
	logger.Errorf("oh crap, an errors occurred: %s", err.Error())
	return
}

// FindAll get a list of links from database.
func (s service) FindAll(ctx context.Context, link findAllRequest) (l []entity.Link, err error) {
	logger := s.logger.With(ctx, "user_id", link.UserID)
	logger.Infof("find links with offset '%d' and limit '%d'", link.Offset, link.Limit)

	err = s.db.Model(&entity.Link{}).Where("user_id = ?", link.UserID).Offset(link.Offset).Limit(link.Limit).Find(&l).Error
	if err == gorm.ErrRecordNotFound {
		logger.Infof("the links with '%d' offset and '%d' limit not found from user_id '%s'", link.Offset, link.Limit, link.UserID)
		return l, e.ErrLinkNotFound
	} else if err == nil {
		return
	}
	logger.Errorf("oh crap, an errors occurred: %s", err.Error())
	return
}

// Update change specific link by ID.
func (s service) Update(ctx context.Context, link updateRequest) (l entity.Link, err error) {
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
		l.URLShort = link.URLShort
		l.URLFull = link.URLFull
		l.Title = link.Title
		l.Active = link.Active
		return
	}

	logger.Errorf("oh crap, an errors occurred: %s", err.Error())
	return
}

// Delete delete a link by ID.
func (s service) Delete(ctx context.Context, link deleteRequest) (err error) {
	logger := s.logger.With(ctx, "user_id", link.UserID)
	logger.Infof("deleting link with id '%s'", link.ID)

	err = s.db.Debug().
		Model(&entity.Link{}).
		Clauses(clause.Returning{}).
		Where("id = ? AND user_id = ?", link.ID, link.UserID).
		Delete(&link).Error

	if err == gorm.ErrRecordNotFound || len(link.URLShort) == 0 {
		logger.Infof("the link id '%s' not found from user_id '%s'", link.ID, link.UserID)
		return e.ErrLinkNotFound
	} else if err == nil {
		return
	}
	logger.Errorf("oh crap, an errors occurred: %s", err.Error())
	return
}
