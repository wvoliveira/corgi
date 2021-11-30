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
	"time"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	AddLink(ctx context.Context, link addLinkRequest) (entity.Link, error)
	FindLinkByID(ctx context.Context, link findLinkByIDRequest) (entity.Link, error)
	FindLinks(ctx context.Context, link findLinksRequest) ([]entity.Link, error)
	UpdateLink(ctx context.Context, link updateLinkRequest) (entity.Link, error)
	DeleteLink(ctx context.Context, link deleteLinkRequest) error

	HTTPAddLink(c *gin.Context)
	HTTPFindLinkByID(c *gin.Context)
	HTTPFindLinks(c *gin.Context)
	HTTPUpdateLink(c *gin.Context)
	HTTPDeleteLink(c *gin.Context)
}

type service struct {
	logger          log.Logger
	db              *gorm.DB
}

// NewService creates a new authentication service.
func NewService(logger log.Logger, db *gorm.DB) Service {
	return service{logger, db}
}

// AddLink create a new shortener link.
func (s service) AddLink(ctx context.Context, link addLinkRequest) (l entity.Link, err error) {
	logger := s.logger.With(ctx, "method", "AddLink", "user_id", link.UserID)

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

// FindLinkByID get a shortener link from ID.
func (s service) FindLinkByID(ctx context.Context, link findLinkByIDRequest) (l entity.Link, err error) {
	logger := s.logger.With(ctx, "method", "FindLinkByID", "user_id", link.UserID)

	err = s.db.Model(&entity.Link{}).Where("link_id = ? AND user_id = ?", link.ID, link.UserID).First(&l).Error
	if err == gorm.ErrRecordNotFound {
		logger.Infof("the link id '%s' not found from user_id '%s'", link.ID, link.UserID)
		return l, e.ErrLinkIDNotFound
	} else if err == nil {
		return
	}
	logger.Errorf("oh crap, an errors occurred: %s", err.Error())
	return
}

// FindLinks get a list of links from database.
func (s service) FindLinks(ctx context.Context, link findLinksRequest) (l []entity.Link, err error) {
	logger := s.logger.With(ctx, "method", "FindLinks", "user_id", link.UserID)

	err = s.db.Model(&entity.Link{}).Where("user_id = ?", link.UserID).Offset(link.Offset).Limit(link.Limit).Find(&l).Error
	if err == gorm.ErrRecordNotFound {
		logger.Infof("the links with '%d' offset and '%d' limit not found from user_id '%s'", link.Offset, link.Limit, link.UserID)
		return l, e.ErrLinkIDNotFound
	} else if err == nil {
		return
	}
	logger.Errorf("oh crap, an errors occurred: %s", err.Error())
	return
}

// UpdateLink update specific link by ID.
func (s service) UpdateLink(ctx context.Context, link updateLinkRequest) (l entity.Link, err error) {
	logger := s.logger.With(ctx, "method", "FindLinks", "user_id", link.UserID)

	err = s.db.Model(&entity.Link{}).Where("link_id = ? AND user_id = ?", link.ID, link.UserID).Updates(&l).Error
	if err == gorm.ErrRecordNotFound {
		logger.Infof("the link id '%s' not found from user_id '%s'", link.ID, link.UserID)
		return l, e.ErrLinkIDNotFound
	} else if err == nil {
		return
	}
	logger.Errorf("oh crap, an errors occurred: %s", err.Error())
	return
}

// DeleteLink delete a link by ID.
func (s service) DeleteLink(ctx context.Context, link deleteLinkRequest) (err error) {
	logger := s.logger.With(ctx, "method", "FindLinks", "user_id", link.UserID)

	err = s.db.Debug().Model(&entity.Link{}).Where("link_id = ? AND user_id = ?", link.ID, link.UserID).Delete(&link).Error
	if err == gorm.ErrRecordNotFound {
		logger.Infof("the link id '%s' not found from user_id '%s'", link.ID, link.UserID)
		return e.ErrLinkIDNotFound
	} else if err == nil {
		return
	}
	logger.Errorf("oh crap, an errors occurred: %s", err.Error())
	return
}
