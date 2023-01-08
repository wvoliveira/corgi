package link

import (
	"errors"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/teris-io/shortid"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	"github.com/wvoliveira/corgi/internal/pkg/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Add(*gin.Context, model.Link) (model.Link, error)
	FindByID(*gin.Context, string, string) (model.Link, error)
	FindAll(*gin.Context, int, int, string, string, string) (int64, int, []model.Link, error)
	Update(*gin.Context, model.Link) (model.Link, error)
	Delete(*gin.Context, string, string) (err error)

	NewHTTP(*gin.RouterGroup)
	HTTPAdd(*gin.Context)
	HTTPFindByID(*gin.Context)
	HTTPFindAll(*gin.Context)
	HTTPUpdate(*gin.Context)
	HTTPDelete(*gin.Context)
}

type service struct {
	db *gorm.DB
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB) Service {
	return service{db}
}

// Add create a new shortener link.
func (s service) Add(c *gin.Context, link model.Link) (m model.Link, err error) {
	l := logger.Logger(c)

	if err = checkLink(link); err != nil {
		return
	}

	// If user is anonymous, create a random ID and blank another fields.
	if link.UserID == "anonymous" {
		sid, _ := shortid.New(1, shortid.DefaultABC, 2342)
		link.Keyword, _ = sid.Generate()
	}

	if link.UserID != "anonymous" {
		if link.Keyword == "" {
			sid, _ := shortid.New(1, shortid.DefaultABC, 2342)
			link.Keyword, _ = sid.Generate()
		}
	}

	err = s.db.Model(&model.Link{}).
		Where("domain = ? AND keyword = ?", link.Domain, link.Keyword).
		Take(&m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {

		m.ID = uuid.New().String()
		m.CreatedAt = time.Now()
		m.Domain = link.Domain
		m.Keyword = link.Keyword
		m.URL = link.URL
		m.Title = link.Title
		m.Active = "true"
		m.UserID = link.UserID

		err = s.db.
			Model(&model.Link{}).
			Create(&m).Error

		if err == nil {
			l.Info().Caller().Msg("short link created with successfully")
		}

		return m, err
	}

	if err == nil {
		l.Warn().Caller().Msg("domain with keyword already exists")
		return m, e.ErrAlreadyExists
	}

	return m, e.ErrInternalServerError
}

// FindByID get a shortener link from ID.
func (s service) FindByID(c *gin.Context, linkID, userID string) (li model.Link, err error) {
	log := logger.Logger(c)

	err = s.db.Model(&model.Link{}).
		Where("id = ? AND user_id = ?", linkID, userID).
		Take(&li).Error

	if err == gorm.ErrRecordNotFound {
		log.Warn().Caller().Msg("link not found")
		return li, e.ErrLinkNotFound

	}

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	return
}

// FindAll get a list of links from database.
func (s service) FindAll(c *gin.Context, offset, limit int, sort, userID, shortenedURL string) (total int64, pages int, links []model.Link, err error) {
	log := logger.Logger(c)

	domain, keyword := util.SplitURL(shortenedURL)

	query := s.db.Model(&model.Link{}).Where("user_id = ?", userID)

	if domain != "" && keyword != "" {
		query = query.Where("domain = ? AND keyword = ?", domain, keyword)
	}

	err = query.
		Count(&total).
		Offset(offset).
		Limit(limit).
		Order(sort).
		Find(&links).Error

	if err == gorm.ErrRecordNotFound {
		log.Info().Caller().Msg(err.Error())
		return total, pages, links, e.ErrLinkNotFound
	}

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return total, pages, links, e.ErrInternalServerError
	}

	pages = int(math.Ceil(float64(total) / float64(limit)))

	return
}

// Update change specific link by ID.
func (s service) Update(c *gin.Context, link model.Link) (m model.Link, err error) {
	log := logger.Logger(c)

	err = s.db.Model(&model.Link{}).
		Where("id = ? AND user_id = ?", link.ID, link.UserID).
		Updates(&link).
		First(&m).Error

	if err == gorm.ErrRecordNotFound {
		log.Info().Caller().Msg(err.Error())
		return m, e.ErrLinkNotFound

	}

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return m, e.ErrInternalServerError
	}

	return
}

// Delete delete a link by ID.
func (s service) Delete(c *gin.Context, linkID, userID string) (err error) {
	log := logger.Logger(c)

	err = s.db.
		Model(&model.Link{}).
		Clauses(clause.Returning{}).
		Where("id = ? AND user_id = ?", linkID, userID).
		Delete(&model.Link{ID: linkID, UserID: userID}).Error

	if err == gorm.ErrRecordNotFound {
		log.Info().Caller().Msg(err.Error())
		return e.ErrLinkNotFound
	}

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return e.ErrInternalServerError
	}

	return
}
