package pwd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Service is a simple CRUD interface for Pwd struct.
type Service interface {
	SignInPwd(ctx context.Context, u Pwd) error
	SignUpPwd(ctx context.Context, u Pwd) error
}

// Pwd represents a single struct for Pwd.
// ID should be globally unique.
type Pwd struct {
	ID        string    `json:"id" gorm:"primaryKey;" example:"eed7df28-5a16-46f0-b5bf-c26071a42ade"`
	CreatedAt time.Time `json:"created_at,omitempty" example:"2021-10-18T00:45:07.818344164-03:00"`
	UpdatedAt time.Time `json:"updated_at,omitempty" example:"2021-10-18T00:49:06.160059334-03:00"`
	LastLogin time.Time `json:"last_login,omitempty" example:"2021-10-20T00:50:00.100059334-03:00"`

	Email    string `json:"email" gorm:"index" example:"oliveira@live.it"`
	Password string `json:"password"`
}

//nolint
var (
	ErrInconsistentIDs     = errors.New("inconsistent IDs")
	ErrAlreadyExists       = errors.New("e-mail already exists")
	ErrNotFound            = errors.New("not found")
	ErrFieldsRequired      = errors.New("fields required: email and password")
	ErrInternalServerError = errors.New("Internal server error")
	ErrUnauthorized        = errors.New("Unauthorized")
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

func (s *dbService) SignInPwd(ctx context.Context, p Pwd) error {
	storedPwd := Pwd{}
	cacheKey := fmt.Sprintf("pwd_email:%s", p.Email)

	if p.Email == "" || p.Password == "" {
		return ErrFieldsRequired
	}

	// check cache
	foo, found := s.c.Get(cacheKey)
	if found {
		storedPwd = foo.(Pwd)
	}

	if !found {
		result := s.db.Model(&storedPwd).Limit(1).Where("email=?", p.Email).Find(&storedPwd)
		if result.RowsAffected == 0 {
			return ErrUnauthorized
		}
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err := bcrypt.CompareHashAndPassword([]byte(storedPwd.Password), []byte(p.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		return ErrUnauthorized
	}

	// store new pwd in in memory cache
	// cacheKey := fmt.Sprintf("pwd_id:%s", p.ID)
	// s.c.Set(cacheKey, p, cache.DefaultExpiration)
	return nil
}

func (s *dbService) SignUpPwd(ctx context.Context, p Pwd) error {
	if p.Email == "" || p.Password == "" {
		return ErrFieldsRequired
	}

	result := s.db.Model(&p).Limit(1).Where("email=?", p.Email).Find(&p)
	if result.RowsAffected > 0 {
		return ErrAlreadyExists // POST = create, don't overwrite
	}

	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p.Password), 8)
	if err != nil {
		return ErrInternalServerError
	}

	p.ID = uuid.New().String()
	p.Password = string(hashedPassword)

	err = s.db.Model(&p).Create(&p).Error
	if err != nil {
		return err
	}

	// store new pwd in in memory cache
	cacheKey := fmt.Sprintf("pwd_email:%s", p.Email)
	s.c.Set(cacheKey, p, cache.DefaultExpiration)
	return nil
}
