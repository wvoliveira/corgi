package app

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service is a simple CRUD interface for user profiles.
type Service interface {
	PostProfile(ctx context.Context, p Profile) error
	GetProfile(ctx context.Context, id string) (Profile, error)
	GetProfiles(ctx context.Context, offset, pageSize int) ([]Profile, error)
	PutProfile(ctx context.Context, id string, p Profile) error
	PatchProfile(ctx context.Context, id string, p Profile) error
	DeleteProfile(ctx context.Context, id string) error
}

// Profile represents a single user profile.
// ID should be globally unique.
type Profile struct {
	ID        string `gorm:"primaryKey;"`
	Name      string `json:"name,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Active    bool `json:"active" gorm:"type:bool;default:true"`
}

func (p *Profile) BeforeCreate(db *gorm.DB) error {
	p.ID = uuid.New().String()
	return nil
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

type dbService struct {
	db *gorm.DB
}

// NewDBService create a new service with gorm DB
func NewDBService(db *gorm.DB) Service {
	return &dbService{
		db: db,
	}
}

func (s *dbService) PostProfile(ctx context.Context, p Profile) error {
	result := s.db.Limit(1).Where("id=?", p.ID).Find(&p)
	if result.RowsAffected > 0 {
		return ErrAlreadyExists // POST = create, don't overwrite
	}

	p.ID = uuid.New().String()
	s.db.Create(&p)
	return nil
}

func (s *dbService) GetProfile(ctx context.Context, id string) (Profile, error) {
	p := Profile{}
	result := s.db.First(&p, id)
	if result.RowsAffected == 0 {
		return p, ErrNotFound
	}
	return p, nil
}

func (s *dbService) GetProfiles(ctx context.Context, offset, pageSize int) ([]Profile, error) {
	var p []Profile
	result := s.db.Offset(offset).Limit(pageSize).Find(&p)
	if result.Error != nil {
		return p, result.Error
	}
	return p, nil
}

func (s *dbService) PutProfile(ctx context.Context, id string, p Profile) error {
	if id != p.ID {
		return ErrInconsistentIDs
	}

	var result *gorm.DB
	if result = s.db.Model(&p).Where("id = ?", id).Updates(&p); result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return s.db.Create(&p).Error
	}

	return nil
}

func (s *dbService) PatchProfile(ctx context.Context, id string, p Profile) error {
	var ep Profile
	result := s.db.First(&ep, id)

	if result.Error != nil {
		return ErrNotFound // PATCH = update existing, don't create
	}

	if id == "" && id != ep.ID {
		return ErrInconsistentIDs
	}

	// We assume that it's not possible to PATCH the ID, and that it's not
	// possible to PATCH any field to its zero value. That is, the zero value
	// means not specified. The way around this is to use e.g. Name *string in
	// the Profile definition. But since this is just a demonstrative example,
	// I'm leaving that out.

	if p.Name != "" {
		ep.Name = p.Name
	}

	if result = s.db.Updates(&ep); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *dbService) DeleteProfile(ctx context.Context, id string) error {
	p := Profile{}

	if result := s.db.First(&p, id); result.RowsAffected == 0 {
		return ErrNotFound
	}

	if result := s.db.Delete(&p, id); result.Error != nil {
		return result.Error
	}
	return nil
}
