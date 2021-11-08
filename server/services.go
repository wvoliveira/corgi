package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Service is a simple CRUD interface for Account struct.
type Service interface {
	SignIn(ctx context.Context, u Account) (Account, error)
	SignUp(ctx context.Context, u Account) error

	AddAccount(ctx context.Context, u Account) (Account, error)
	FindAccountByID(ctx context.Context, id string) (Account, error)
	FindAccounts(ctx context.Context, offset, pageSize int) ([]Account, error)
	UpdateOrCreateAccount(ctx context.Context, id string, u Account) error
	UpdateAccount(ctx context.Context, id string, u Account) error
	DeleteAccount(ctx context.Context, id string) error

	AddURL(ctx context.Context, u URL) (URL, error)
	FindURLByID(ctx context.Context, id string) (URL, error)
	FindURLs(ctx context.Context, offset, pageSize int) ([]URL, error)
	UpdateOrCreateURL(ctx context.Context, id string, u URL) error
	UpdateURL(ctx context.Context, id string, u URL) error
	DeleteURL(ctx context.Context, id string) error
}

type service struct {
	db    *gorm.DB
	cache *cache.Cache
}

// NewService create a new service with database and cache.
func NewService(db *gorm.DB, c *cache.Cache) Service {
	return &service{
		db:    db,
		cache: c,
	}
}

func (s *service) SignIn(ctx context.Context, p Account) (Account, error) {
	storedAccount := Account{}
	cacheKey := fmt.Sprintf("pwd_email:%s", p.Email)

	if p.Email == "" || p.Password == "" {
		return p, ErrFieldsRequired
	}

	// check cache
	foo, found := s.cache.Get(cacheKey)
	if found {
		storedAccount = foo.(Account)
	}

	if !found {
		result := s.db.Model(&storedAccount).Limit(1).Where("email=?", p.Email).Find(&storedAccount)
		if result.RowsAffected == 0 {
			return p, ErrUnauthorized
		}
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err := bcrypt.CompareHashAndPassword([]byte(storedAccount.Password), []byte(p.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		return p, ErrUnauthorized
	}

	// Get request from context.
	r := ctx.Value(ctxRequestKey{}).(*http.Request)

	// Get session from request.
	session, _ := store.Get(r, "session")
	session.Values["account_id"] = storedAccount.ID

	// Set session to Account struct session.
	p.Session = session

	cacheSessionKey := fmt.Sprintf("pwd_session_account:%s", p.ID)
	s.cache.Set(cacheSessionKey, session, cache.DefaultExpiration)

	// set session token to Account object
	p.Session = session

	// store new pwd in in memory cache
	// cacheKey := fmt.Sprintf("pwd_id:%s", p.ID)
	// s.cache.Set(cacheKey, p, cache.DefaultExpiration)
	return p, nil
}

func (s *service) SignUp(ctx context.Context, p Account) error {
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
	s.cache.Set(cacheKey, p, cache.DefaultExpiration)
	return nil
}

/*
	Create a new account.
*/
func (s *service) AddAccount(ctx context.Context, p Account) (Account, error) {
	if p.Email == "" || p.Password == "" {
		return p, errors.New("fields required: email and password")
	}

	result := s.db.Model(&p).Limit(1).Where("email=?", p.Email).Find(&p)
	if result.RowsAffected > 0 {
		return p, ErrAlreadyExists // POST = create, don't overwrite
	}

	p.ID = uuid.New().String()
	err := s.db.Model(&p).Create(&p).Error
	if err != nil {
		return p, err
	}

	// store new account in in memory cache
	cacheKey := fmt.Sprintf("account_id:%s", p.ID)
	s.cache.Set(cacheKey, p, cache.DefaultExpiration)
	return p, nil
}

/*
	Get account by ID.
*/
func (s *service) FindAccountByID(ctx context.Context, id string) (Account, error) {
	u := Account{}
	cacheKey := fmt.Sprintf("account_id:%s", id)

	// check cache
	foo, found := s.cache.Get(cacheKey)
	if found {
		return foo.(Account), nil
	}

	// get in db
	result := s.db.Model(&u).Where("id = ?", id).First(&u)
	if result.RowsAffected == 0 {
		return u, ErrNotFound
	}

	// if found, set cache
	s.cache.Set(cacheKey, u, cache.DefaultExpiration)
	return u, nil
}

/*
	Get a list of accounts.
*/
func (s *service) FindAccounts(ctx context.Context, offset, pageSize int) ([]Account, error) {
	var u []Account

	// get account in db
	result := s.db.Model(&u).Offset(offset).Limit(pageSize).Find(&u)
	if result.Error != nil {
		return u, result.Error
	}
	return u, nil
}

/*
	Update or create a new account.
*/
func (s *service) UpdateOrCreateAccount(ctx context.Context, id string, reqAccount Account) error {
	var dbAccount Account
	var result *gorm.DB

	cacheKey := fmt.Sprintf("account_id:%s", id)

	if _, err := uuid.Parse(id); err != nil {
		return ErrInconsistentIDs
	}

	result = s.db.Model(&reqAccount).Where("id = ?", id).First(&dbAccount)

	if result.RowsAffected == 0 {
		if reqAccount.Email == "" {
			return errors.New("fields required: email and password")
		}

		if result = s.db.Model(&reqAccount).Create(&reqAccount); result.Error != nil {
			return result.Error
		}
	}

	if result = s.db.Model(&dbAccount).Where("id=?", id).Save(&dbAccount); result.Error != nil {
		return result.Error
	}

	s.cache.Set(cacheKey, dbAccount, cache.DefaultExpiration)
	return nil
}

/*
	Updates a account that already exists. Do not create.
*/
func (s *service) UpdateAccount(ctx context.Context, id string, reqAccount Account) error {
	var dbAccount Account

	result := s.db.Model(&reqAccount).First(&dbAccount, id)
	cacheKey := fmt.Sprintf("account_id:%s", id)

	if result.Error != nil {
		return ErrNotFound
	}

	if result = s.db.Model(&reqAccount).Updates(&reqAccount); result.Error != nil {
		return result.Error
	}

	s.cache.Set(cacheKey, reqAccount, cache.DefaultExpiration)
	return nil
}

/*
	Delete account by ID.
*/
func (s *service) DeleteAccount(ctx context.Context, id string) error {
	u := Account{}
	cacheKey := fmt.Sprintf("account_id:%s", id)

	if result := s.db.Model(&u).Where("id = ?", id).First(&u); result.RowsAffected == 0 {
		return ErrNotFound
	}

	if result := s.db.Model(&u).Where("id = ?", id).Delete(&u); result.Error != nil {
		return result.Error
	}

	// Delete URL in cache.
	s.cache.Delete(cacheKey)
	return nil
}

func (s *service) AddURL(ctx context.Context, u URL) (URL, error) {
	if u.Keyword == "" || u.URL == "" || u.Title == "" || u.OwnerID == "" {
		return URL{}, errors.New("fields required: keyword, url, title and owner_id")
	}

	result := s.db.Model(&u).Limit(1).Where("keyword=?", u.Keyword).Find(&u)
	if result.RowsAffected > 0 {
		return URL{}, ErrAlreadyExists // POST = create, don't overwrite
	}

	u.ID = uuid.New().String()
	err := s.db.Model(&u).Create(&u).Error
	if err != nil {
		return URL{}, err
	}

	// Store new URL in memory cache.
	cacheKey := fmt.Sprintf("url_id:%s", u.ID)
	s.cache.Set(cacheKey, u, cache.DefaultExpiration)
	return u, nil
}

func (s *service) FindURLByID(ctx context.Context, id string) (URL, error) {
	u := URL{}
	cacheKey := fmt.Sprintf("url_id:%s", id)

	// Check if URL exists in cache.
	foo, found := s.cache.Get(cacheKey)
	if found {
		return foo.(URL), nil
	}

	// Get URL in database.
	result := s.db.Model(&u).Where("id = ?", id).First(&u)
	if result.RowsAffected == 0 {
		return u, ErrNotFound
	}

	// If found, set cache.
	s.cache.Set(cacheKey, u, cache.DefaultExpiration)
	return u, nil
}

func (s *service) FindURLs(ctx context.Context, offset, pageSize int) ([]URL, error) {
	var u []URL

	// get url in db
	result := s.db.Model(&u).Offset(offset).Limit(pageSize).Find(&u)
	if result.Error != nil {
		return u, result.Error
	}
	return u, nil
}

/*
	Update or create a new URL.
*/
func (s *service) UpdateOrCreateURL(ctx context.Context, id string, reqURL URL) error {
	var dbURL URL
	var result *gorm.DB
	cacheKey := fmt.Sprintf("url_id:%s", id)

	_, err := uuid.Parse(id)

	if err != nil {
		return ErrInconsistentIDs
	}

	result = s.db.Model(&reqURL).Where("id = ?", id).First(&dbURL)

	if result.RowsAffected == 0 {
		if reqURL.Keyword == "" || reqURL.URL == "" || reqURL.Title == "" || reqURL.OwnerID == "" {
			return errors.New("fields required: keyword, url, title and owner_id")
		}

		if result = s.db.Model(&reqURL).Create(&reqURL); result.Error != nil {
			return result.Error
		}
	}

	if result = s.db.Model(&dbURL).Where("id = ?", id).Save(&dbURL); result.Error != nil {
		return result.Error
	}

	s.cache.Set(cacheKey, dbURL, cache.DefaultExpiration)
	return nil
}

/*
	Updates a URL that already exists. Do not create.
*/
func (s *service) UpdateURL(ctx context.Context, id string, reqURL URL) error {
	var dbURL URL

	result := s.db.Model(&reqURL).First(&dbURL, id)
	cacheKey := fmt.Sprintf("url_id:%s", id)

	if result.Error != nil {
		return ErrNotFound
	}

	if id == "" && id != dbURL.ID {
		return ErrInconsistentIDs
	}

	if result = s.db.Model(&dbURL).Updates(&dbURL); result.Error != nil {
		return result.Error
	}

	s.cache.Set(cacheKey, dbURL, cache.DefaultExpiration)
	return nil
}

/*
	Delete a URL by ID.
*/
func (s *service) DeleteURL(ctx context.Context, id string) error {
	u := URL{}
	cacheKey := fmt.Sprintf("url_id:%s", id)

	if result := s.db.Model(&u).Where("id = ?", id).First(&u); result.RowsAffected == 0 {
		return ErrNotFound
	}

	if result := s.db.Model(&u).Where("id = ?", id).Delete(&u); result.Error != nil {
		return result.Error
	}

	s.cache.Delete(cacheKey)
	return nil
}
