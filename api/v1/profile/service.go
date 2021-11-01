package profile

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

// Service is a simple CRUD interface for Profile struct.
type Service interface {
	PostProfile(ctx context.Context, u Profile) error
	GetProfile(ctx context.Context, id string) (Profile, error)
	GetProfiles(ctx context.Context, offset, pageSize int) ([]Profile, error)
	PutProfile(ctx context.Context, id string, u Profile) error
	PatchProfile(ctx context.Context, id string, u Profile) error
	DeleteProfile(ctx context.Context, id string) error
}

// Profile represents a single struct for Profile.
// ID should be globally unique.
type Profile struct {
	ID        string    `json:"id" gorm:"primaryKey;" example:"eed7df28-5a16-46f0-b5bf-c26071a42ade"`
	CreatedAt time.Time `json:"created_at,omitempty" example:"2021-10-18T00:45:07.818344164-03:00"`
	UpdatedAt time.Time `json:"updated_at,omitempty" example:"2021-10-18T00:49:06.160059334-03:00"`
	LastLogin time.Time `json:"last_login,omitempty" example:"2021-10-20T00:50:00.100059334-03:00"`

	Name     string   `json:"name" example:"Wellington Oliveira"`
	Email    string   `json:"email" gorm:"index" example:"oliveira@live.it"`
	Password string   `json:"-"`
	Active   *bool    `json:"active" gorm:"type:bool;default:true" example:"false"`
	Roles    []string `json:"roles" gorm:"array"`
	Tags     []string `json:"tags" gorm:"array"`
}

// PostProfile struct when response request
type PostProfile struct {
	ID   string `json:"id"`
	Name string `json:"keyword"`
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

func (s *dbService) PostProfile(ctx context.Context, u Profile) error {
	if u.Name == "" || u.Email == "" || u.Password == "" {
		return errors.New("fields required: name, email and password")
	}

	result := s.db.Model(&u).Limit(1).Where("email=?", u.Email).Find(&u)
	if result.RowsAffected > 0 {
		return ErrAlreadyExists // POST = create, don't overwrite
	}

	u.ID = uuid.New().String()
	err := s.db.Model(&u).Create(&u).Error
	if err != nil {
		return err
	}

	// store new profile in in memory cache
	cacheKey := fmt.Sprintf("profile_id:%s", u.ID)
	s.c.Set(cacheKey, u, cache.DefaultExpiration)
	return nil
}

func (s *dbService) GetProfile(ctx context.Context, id string) (Profile, error) {
	u := Profile{}
	cacheKey := fmt.Sprintf("profile_id:%s", id)

	// check cache
	foo, found := s.c.Get(cacheKey)
	if found {
		return foo.(Profile), nil
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

func (s *dbService) GetProfiles(ctx context.Context, offset, pageSize int) ([]Profile, error) {
	var u []Profile

	// get profile in db
	result := s.db.Model(&u).Offset(offset).Limit(pageSize).Find(&u)
	if result.Error != nil {
		return u, result.Error
	}
	return u, nil
}

func (s *dbService) PutProfile(ctx context.Context, id string, u Profile) error {
	// PUT = Update or create
	var eu Profile
	var result *gorm.DB
	cacheKey := fmt.Sprintf("profile_id:%s", id)

	_, err := uuid.Parse(id)

	if err != nil {
		return ErrInconsistentIDs
	}

	result = s.db.Model(&u).Where("id = ?", id).First(&eu)

	if result.RowsAffected == 0 {
		if u.Name == "" || u.Email == "" || u.Password == "" {
			return errors.New("fields required: name, email and password")
		}

		if result = s.db.Model(&u).Create(&u); result.Error != nil {
			return result.Error
		}
	}

	if u.Name != "" {
		eu.Name = u.Name
	}

	if u.Email != "" {
		eu.Email = u.Email
	}

	if u.Password != "" {
		eu.Password = u.Password
	}

	if u.Active != nil {
		eu.Active = u.Active
	}

	if result = s.db.Model(&eu).Where("id=?", id).Save(&eu); result.Error != nil {
		return result.Error
	}

	// set cache
	s.c.Set(cacheKey, eu, cache.DefaultExpiration)
	return nil
}

func (s *dbService) PatchProfile(ctx context.Context, id string, u Profile) error {
	// PATCH = update existing, don't create
	var eu Profile
	result := s.db.Model(&u).First(&eu, id)
	cacheKey := fmt.Sprintf("profile_id:%s", id)

	if result.Error != nil {
		return ErrNotFound
	}

	if id == "" && id != eu.ID {
		return ErrInconsistentIDs
	}

	// We assume that it's not possible to PATCH the ID, and that it's not
	// possible to PATCH any field to its zero value. That is, the zero value
	// means not specified. The way around this is to use e.g. Name *string in
	// the Profile definition. But since this is just a demonstrative example,
	// I'm leaving that out.

	if u.Name != "" {
		eu.Name = u.Name
	}

	if u.Email != "" {
		eu.Email = u.Email
	}

	if u.Password != "" {
		eu.Password = u.Password
	}

	if u.Active != nil {
		eu.Active = u.Active
	}

	if result = s.db.Model(&u).Updates(&eu); result.Error != nil {
		return result.Error
	}

	// set in cache
	s.c.Set(cacheKey, eu, cache.DefaultExpiration)
	return nil
}

func (s *dbService) DeleteProfile(ctx context.Context, id string) error {
	u := Profile{}
	cacheKey := fmt.Sprintf("profile_id:%s", id)

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
