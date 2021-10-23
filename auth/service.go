package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

// Service is a simple CRUD interface for Auth struct.
type Service interface {
	PostSignup(ctx context.Context, a Auth) error
	PostLogin(ctx context.Context, a Auth) error
	PostLogout(ctx context.Context, a Auth) error
	PostRefresh(ctx context.Context, a Auth) error
}

// Auth represents a single struct for Auth.
// ID should be globally unique.
type Auth struct {
	ID       string `json:"id" gorm:"primaryKey;unique" example:"eed7df28-5a16-46f0-b5bf-c26071a42ade"`
	Email    string `json:"email" gorm:"index" example:"oliveira@live.it"`
	Password string `json:"password" gorm:"not null;"`
}

type AccessDetails struct {
	AccessUuid string
	UserId     int64
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AuthPost struct {
	Email string `json:"email" example:"oliveira@live.it"`
}

//nolint
var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("user already exists")
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

func (s *dbService) PostSignup(ctx context.Context, a Auth) error {
	var ea Auth

	if a.Email == "" || a.Password == "" {
		return errors.New("Fields required: email and password")
	}

	result := s.db.Limit(1).Where("email=?", a.Email).Find(&ea)
	if result.RowsAffected > 0 {
		return ErrAlreadyExists
	}

	a.ID = uuid.New().String()
	err := s.db.Create(&a).Error
	if err != nil {
		return err
	}

	// store new auth in in memory cache
	cacheKey := fmt.Sprintf("auth_id:%s", a.ID)
	s.c.Set(cacheKey, a, cache.DefaultExpiration)
	return nil
}

func (s *dbService) PostLogin(ctx context.Context, a Auth) error {
	var ea Auth

	if a.Email == "" || a.Password == "" {
		return errors.New("Fields required: email and password")
	}

	result := s.db.Limit(1).Where("email=?", u.Email).Find(&ea)
	if result.RowsAffected == 0 {
		return errors.New("User doesnt exist!")
	}

	err := bcrypt.CompareHashAndPassword([]byte(ea.Password), []byte(a.Password))
	if err != nil {
		return errors.New("Password doesnt match!")
	}
	return nil

	// store new user in in memory cache
	// cacheKey := fmt.Sprintf("user_id:%s", u.ID)
	// s.c.Set(cacheKey, u, cache.DefaultExpiration)
	// return nil
}

func (s *dbService) PostLogout(ctx context.Context, a Auth) error  { return nil }
func (s *dbService) PostRefresh(ctx context.Context, a Auth) error { return nil }
