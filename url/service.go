package url

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

// Service is a simple CRUD interface for URL struct.
type Service interface {
	PostURL(ctx context.Context, u URL) error
	GetURL(ctx context.Context, id string) (URL, error)
	GetURLs(ctx context.Context, offset, pageSize int) ([]URL, error)
	PutURL(ctx context.Context, id string, u URL) error
	PatchURL(ctx context.Context, id string, u URL) error
	DeleteURL(ctx context.Context, id string) error
}

// URL represents a single struct for URL.
// ID should be globally unique.
type URL struct {
	ID        string    `json:"id" gorm:"primaryKey;" example:"eed7df28-5a16-46f0-b5bf-c26071a42ade"`
	Keyword   string    `json:"keyword" gorm:"index" example:"google"`
	URL       string    `json:"url" example:"https://www.google.com"`
	Title     string    `json:"title" example:"Google Home"`
	Active    *bool     `json:"active" gorm:"type:bool;default:true" example:"false"`
	OwnerID   string    `json:"owner_id" example:"5ca04a43-ff3c-4154-a8ad-02e2e906a847"`
	CreatedAt time.Time `json:"created_at" example:"2021-10-18T00:45:07.818344164-03:00"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-10-18T00:49:06.160059334-03:00"`
}

type PostURL struct {
	Keyword string `json:"keyword" gorm:"index" example:"google"`
	URL     string `json:"url" example:"https://www.google.com"`
	Title   string `json:"title" example:"Google Home"`
}

//nolint
var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
	ErrFieldsRequired  = errors.New("fields required: ")
)

type dbService struct {
	db *gorm.DB
	c  *cache.Cache
}

// NewDBService create a new service with gorm DB
func NewDBService(db *gorm.DB, c *cache.Cache) Service {
	return &dbService{
		db: db,
		c:  c,
	}
}

func (s *dbService) PostURL(ctx context.Context, u URL) error {
	if u.Keyword == "" || u.URL == "" || u.Title == "" || u.OwnerID == "" {
		return errors.New("fields required: keyword, url, title and owner_id")
	}

	result := s.db.Limit(1).Where("keyword=?", u.Keyword).Find(&u)
	if result.RowsAffected > 0 {
		return ErrAlreadyExists // POST = create, don't overwrite
	}

	u.ID = uuid.New().String()
	err := s.db.Create(&u).Error
	if err != nil {
		return err
	}

	// store new url in in memory cache
	cacheKey := fmt.Sprintf("url_id:%s", u.ID)
	s.c.Set(cacheKey, u, cache.DefaultExpiration)
	return nil
}

func (s *dbService) GetURL(ctx context.Context, id string) (URL, error) {
	u := URL{}
	cacheKey := fmt.Sprintf("url_id:%s", id)

	// check cache
	foo, found := s.c.Get(cacheKey)
	if found {
		return foo.(URL), nil
	}

	// get in db
	result := s.db.Model(&u).Where("id = ?", id).First(&u)
	if result.RowsAffected == 0 {
		return u, ErrNotFound
	}

	// if found, set cache
	s.c.Set(cacheKey, u, cache.DefaultExpiration)
	return u, nil
}

func (s *dbService) GetURLs(ctx context.Context, offset, pageSize int) ([]URL, error) {
	var u []URL

	// get url in db
	result := s.db.Offset(offset).Limit(pageSize).Find(&u)
	if result.Error != nil {
		return u, result.Error
	}
	return u, nil
}

func (s *dbService) PutURL(ctx context.Context, id string, u URL) error {
	// PUT = Update or create
	var eu URL
	var result *gorm.DB
	cacheKey := fmt.Sprintf("url_id:%s", id)

	if !isValidUUID(id) {
		return ErrInconsistentIDs
	}

	result = s.db.Model(&u).Where("id = ?", id).First(&eu)

	if result.RowsAffected == 0 {
		if u.Keyword == "" || u.URL == "" || u.Title == "" || u.OwnerID == "" {
			return errors.New("fields required: keyword, url, title and owner_id")
		}

		if result = s.db.Create(&u); result.Error != nil {
			return result.Error
		}
	}

	if u.Keyword != "" {
		eu.Keyword = u.Keyword
	}

	if u.URL != "" {
		eu.URL = u.URL
	}

	if u.Title != "" {
		eu.Title = u.Title
	}

	if u.Active != nil {
		eu.Active = u.Active
	}

	if result = s.db.Model(&eu).Where("id = ?", id).Save(&eu); result.Error != nil {
		return result.Error
	}

	// set cache
	s.c.Set(cacheKey, eu, cache.DefaultExpiration)
	return nil
}

func (s *dbService) PatchURL(ctx context.Context, id string, u URL) error {
	// PATCH = update existing, don't create
	var eu URL
	result := s.db.First(&eu, id)
	cacheKey := fmt.Sprintf("url_id:%s", id)

	if result.Error != nil {
		return ErrNotFound
	}

	if id == "" && id != eu.ID {
		return ErrInconsistentIDs
	}

	// We assume that it's not possible to PATCH the ID, and that it's not
	// possible to PATCH any field to its zero value. That is, the zero value
	// means not specified. The way around this is to use e.g. Name *string in
	// the URL definition. But since this is just a demonstrative example,
	// I'm leaving that out.

	if u.Keyword != "" {
		eu.Keyword = u.Keyword
	}

	if u.URL != "" {
		eu.URL = u.URL
	}

	if u.Title != "" {
		eu.Title = u.Title
	}

	if u.Active != nil {
		eu.Active = u.Active
	}

	if result = s.db.Updates(&eu); result.Error != nil {
		return result.Error
	}

	// set in cache
	s.c.Set(cacheKey, eu, cache.DefaultExpiration)
	return nil
}

func (s *dbService) DeleteURL(ctx context.Context, id string) error {
	u := URL{}
	cacheKey := fmt.Sprintf("url_id:%s", id)

	if result := s.db.Model(&u).Where("id = ?", id).First(&u); result.RowsAffected == 0 {
		return ErrNotFound
	}

	if result := s.db.Model(&u).Where("id = ?", id).Delete(&u); result.Error != nil {
		return result.Error
	}

	// delete cache
	s.c.Delete(cacheKey)
	return nil
}
