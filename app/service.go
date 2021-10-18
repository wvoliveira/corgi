package app

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
	ID        string    `json:"id" gorm:"primaryKey;"`
	Keyword   string    `json:"keyword" gorm:"index"`
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	Active    *bool     `json:"active" gorm:"type:bool;default:true"`
	OwnerID   string    `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	cache_key := fmt.Sprintf("url_id:%s", u.ID)
	s.c.Set(cache_key, u, cache.DefaultExpiration)
	return nil
}

func (s *dbService) GetURL(ctx context.Context, id string) (URL, error) {
	u := URL{}
	cache_key := fmt.Sprintf("url_id:%s", id)

	// check cache
	foo, found := s.c.Get(cache_key)
	if found {
		return foo.(URL), nil
	}

	// get in db
	result := s.db.Model(&u).Where("id = ?", id).First(&u)
	if result.RowsAffected == 0 {
		return u, ErrNotFound
	}

	// if found, set cache
	s.c.Set(cache_key, u, cache.DefaultExpiration)
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
	var result *gorm.DB
	cache_key := fmt.Sprintf("url_id:%s", id)

	if id != u.ID {
		return ErrInconsistentIDs
	}

	// Fix: create if not exist
	// PUT = create or update
	if result = s.db.Model(&u).Where("id = ?", id).Updates(&u); result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		result = s.db.Create(&u)
	}

	if result.Error != nil {
		return result.Error
	}

	// set cache
	s.c.Set(cache_key, u, cache.DefaultExpiration)
	return nil
}

func (s *dbService) PatchURL(ctx context.Context, id string, u URL) error {
	var eu URL
	result := s.db.First(&eu, id)
	cache_key := fmt.Sprintf("url_id:%s", id)

	if result.Error != nil {
		return ErrNotFound // PATCH = update existing, don't create
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
	s.c.Set(cache_key, eu, cache.DefaultExpiration)
	return nil
}

func (s *dbService) DeleteURL(ctx context.Context, id string) error {
	u := URL{}
	cache_key := fmt.Sprintf("url_id:%s", id)

	if result := s.db.Model(&u).Where("id = ?", id).First(&u); result.RowsAffected == 0 {
		return ErrNotFound
	}

	if result := s.db.Model(&u).Where("id = ?", id).Delete(&u); result.Error != nil {
		return result.Error
	}

	// delete cache
	s.c.Delete(cache_key)
	return nil
}
