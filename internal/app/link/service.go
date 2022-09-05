package link

/*
TODO:
- Create logic to create link anonymously
*/

import (
	"context"
	"errors"
	"math"
	"net/http"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/teris-io/shortid"
	"github.com/wvoliveira/corgi/internal/app/entity"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Add(ctx context.Context, payload entity.Link) (link entity.Link, err error)
	FindByID(ctx context.Context, linkID, userID string) (link entity.Link, err error)
	FindAll(ctx context.Context, offset, limit int, sort, userID string) (total int64, pages int, links []entity.Link, err error)
	Update(ctx context.Context, payload entity.Link) (link entity.Link, err error)
	Delete(ctx context.Context, linkID, userID string) (err error)

	NewHTTP(r *mux.Router)
	HTTPAdd(w http.ResponseWriter, r *http.Request)
	HTTPFindByID(w http.ResponseWriter, r *http.Request)
	HTTPFindAll(w http.ResponseWriter, r *http.Request)
	HTTPUpdate(w http.ResponseWriter, r *http.Request)
	HTTPDelete(w http.ResponseWriter, r *http.Request)
}

type service struct {
	db      *gorm.DB
	secret  string
	store   *sessions.CookieStore
	enforce *casbin.Enforcer
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB, secret string, store *sessions.CookieStore, enforce *casbin.Enforcer) Service {
	return service{db, secret, store, enforce}
}

// Add create a new shortener link.
func (s service) Add(ctx context.Context, link entity.Link) (li entity.Link, err error) {
	l := logger.Logger(ctx)

	if err = checkLink(link); err != nil {
		return
	}

	// If userID is anonymous, create a random ID.
	// And blank another fields.
	// TODO: with anonymous access, keep trying to create a unique keyword.
	if link.UserID == "anonymous" {
		sid, _ := shortid.New(1, shortid.DefaultABC, 2342)
		link.Keyword, _ = sid.Generate()

		// In anonymous, for small database, we check if URL already exists.
		err = s.db.Model(&entity.Link{}).
			Where("domain = ? AND keyword = ? OR (domain = ? AND url = ?)", link.Domain, link.Keyword, link.Domain, link.URL).
			Take(&li).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {

			li.ID = uuid.New().String()
			li.CreatedAt = time.Now()
			li.Domain = link.Domain
			li.Keyword = link.Keyword
			li.URL = link.URL
			li.Active = "true"
			li.UserID = link.UserID

			err = s.db.Model(&entity.Link{}).Create(&li).Error
			return li, err
		} else if err == nil {
			l.Warn().Caller().Msg("domain with keyword already exists")
			return li, e.ErrAnonymousURLAlreadyExists
		}
		l.Error().Caller().Msg(err.Error())
	}

	// With a real user authentication.
	if link.Domain == "" {
		return li, e.ErrLinkInvalidDomain
	}
	if link.Keyword == "" {
		return li, e.ErrLinkInvalidKeyword
	}

	err = s.db.Model(&entity.Link{}).Where("domain = ? AND keyword = ?", link.Domain, link.Keyword).Take(&li).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {

		li.ID = uuid.New().String()
		li.CreatedAt = time.Now()
		li.Domain = link.Domain
		li.Keyword = link.Keyword
		li.Title = link.Title
		li.URL = link.URL
		li.Active = "true"
		li.UserID = link.UserID

		err = s.db.Model(&entity.Link{}).Create(&li).Error
		return li, err
	} else if err == nil {
		l.Warn().Caller().Msg("domain with keyword already exists")
		return li, e.ErrLinkAlreadyExists
	}
	l.Error().Caller().Msg(err.Error())
	return li, e.ErrInternalServerError
}

// FindByID get a shortener link from ID.
func (s service) FindByID(ctx context.Context, linkID, userID string) (li entity.Link, err error) {
	l := logger.Logger(ctx)

	err = s.db.Model(&entity.Link{}).
		Where("id = ? AND user_id = ?", linkID, userID).
		Take(&li).Error

	if err == gorm.ErrRecordNotFound {
		l.Warn().Caller().Msg("link not found")
		return li, e.ErrLinkNotFound
	} else if err == nil {
		return
	}
	l.Error().Caller().Msg(err.Error())
	return
}

// FindAll get a list of links from database.
func (s service) FindAll(ctx context.Context, offset, limit int, sort, userID string) (total int64, pages int, links []entity.Link, err error) {
	l := logger.Logger(ctx)

	err = s.db.Model(&entity.Link{}).Where("user_id = ?", userID).
		Count(&total).
		Offset(offset).
		Limit(limit).
		Order(sort).
		Find(&links).Error

	if err == gorm.ErrRecordNotFound {
		l.Info().Caller().Msg("links not found")
		return total, pages, links, e.ErrLinkNotFound
	} else if err != nil {
		l.Error().Caller().Msg(err.Error())
		return
	}

	pages = int(math.Ceil(float64(total) / float64(limit)))
	return
}

// Update change specific link by ID.
func (s service) Update(ctx context.Context, link entity.Link) (li entity.Link, err error) {
	l := logger.Logger(ctx)

	err = s.db.Model(&entity.Link{}).
		Where("id = ? AND user_id = ?", link.ID, link.UserID).
		Updates(&link).
		First(&li).Error

	if err == gorm.ErrRecordNotFound {
		l.Info().Caller().Msg("link not found")
		return li, e.ErrLinkNotFound
	} else if err != nil {
		return li, e.ErrInternalServerError
	}
	return
}

// Delete delete a link by ID.
func (s service) Delete(ctx context.Context, linkID, userID string) (err error) {
	l := logger.Logger(ctx)

	err = s.db.Debug().
		Model(&entity.Link{}).
		Clauses(clause.Returning{}).
		Where("id = ? AND user_id = ?", linkID, userID).
		Delete(&entity.Link{ID: linkID, UserID: userID}).Error

	if err == gorm.ErrRecordNotFound {
		l.Info().Caller().Msg("link not found")
		return e.ErrLinkNotFound
	} else if err == nil {
		return
	}
	l.Error().Caller().Msg(err.Error())
	return
}
