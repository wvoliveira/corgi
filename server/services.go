package server

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Service store all methods. Yeah monolithic.
type Service struct {
	db     *gorm.DB
	cache  *cache.Cache
	secret string
}

// NewService create a new service with database and cache.
func NewService(secretKey string, db *gorm.DB, c *cache.Cache) Service {
	return Service{
		db:     db,
		cache:  c,
		secret: secretKey,
	}
}

// SignIn login with email and password.
func (s Service) SignIn(p Account) (Account, error) {
	sa := Account{}
	cacheKey := fmt.Sprintf("pwd_email:%s", p.Email)

	// Check if e-mail and password was sended.
	if p.Email == "" || p.Password == "" {
		return p, ErrFieldsRequired
	}

	// Check memory cache.
	foo, found := s.cache.Get(cacheKey)
	if found {
		sa = foo.(Account)
	}

	// If not found in cache memory, search in database (more slowly).
	if !found {
		result := s.db.Model(&sa).Limit(1).Where("email=?", p.Email).Find(&sa)
		if result.RowsAffected == 0 {
			return p, ErrUnauthorized
		}
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err := bcrypt.CompareHashAndPassword([]byte(sa.Password), []byte(p.Password)); err != nil {
		return p, ErrUnauthorized
	}

	// Generate JWT with specific claims.
	tokenHash, err := generateJWT(s.secret, sa)
	if err != nil {
		return sa, err
	}

	// Set token account with JWT.
	sa.Token = tokenHash

	// Store the account in cache.
	s.cache.Set(cacheKey, sa, cache.DefaultExpiration)

	return sa, nil
}

// SignUp register with e-mail and password.
func (s Service) SignUp(p Account) error {
	// Check if e-mail and password was sended.
	// TODO: check e-mail pattern.
	if p.Email == "" || p.Password == "" {
		return ErrFieldsRequired
	}

	// Key for save account in cache.
	cacheKey := fmt.Sprintf("account_email:%s", p.Email)

	// Check if account is into the cache.
	_, found := s.cache.Get(cacheKey)
	if found {
		return ErrAlreadyExists
	}

	// Yeah, I don't using 'else' statement.
	// Dont overwrite account.
	if !found {
		result := s.db.Model(&p).Limit(1).Where("email=?", p.Email).Find(&p)
		if result.RowsAffected > 0 {
			return ErrAlreadyExists
		}
	}

	// Salt and hash the password using the bcrypt algorithm.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p.Password), 8)
	if err != nil {
		return ErrInternalServerError
	}

	// Create a random ID and default role for new user.
	p.ID = uuid.New().String()
	p.Role = "user"
	p.Password = string(hashedPassword)

	// Create a new account.
	err = s.db.Model(&p).Create(&p).Error
	if err != nil {
		return err
	}

	// Store the new account in cache.
	s.cache.Set(cacheKey, p, cache.DefaultExpiration)
	return nil
}

// AddAccount create a new account.
func (s Service) AddAccount(auth Account, payload Account) (Account, error) {
	// Only admin can create a new account without signup process.
	// TODO: create more roles without hardcoded.
	if auth.Role != "admin" {
		return auth, ErrOnlyAdmin
	}

	// Check fields in payload.
	if payload.Email == "" || payload.Password == "" {
		return payload, ErrFieldsRequired
	}

	// Check if account exists. Don't overwrite.
	result := s.db.Model(&payload).Limit(1).Where("email=?", payload.Email).Find(&payload)
	if result.RowsAffected > 0 {
		return payload, ErrAlreadyExists
	}

	// Set a random ID and create account.
	payload.ID = uuid.New().String()
	err := s.db.Model(&payload).Create(&payload).Error
	if err != nil {
		return payload, err
	}

	// Store new account in memory cache.
	cacheKey := fmt.Sprintf("account_id:%s", payload.ID)
	s.cache.Set(cacheKey, payload, cache.DefaultExpiration)
	return payload, nil
}

// FindAccountByID find a account with specific ID.
func (s Service) FindAccountByID(auth Account, id string) (acc Account, err error) {
	// Only admin can view another accounts.
	// TODO: create more roles without hardcoded.
	if auth.Role != "admin" && auth.ID != id {
		return acc, ErrOnlyAdmin
	}

	cacheKey := fmt.Sprintf("account_id:%s", id)

	// Check cache if account is present.
	foo, found := s.cache.Get(cacheKey)
	if found {
		return foo.(Account), err
	}

	// Get account from databse.
	result := s.db.Model(&acc).Where("id = ?", id).First(&acc)
	if result.RowsAffected == 0 {
		return acc, ErrNotFound
	}

	// If found, send do memory cache.
	s.cache.Set(cacheKey, acc, cache.DefaultExpiration)
	return
}

// FindAccounts Get a list of accounts.
func (s Service) FindAccounts(auth Account, offset, pageSize int) (accs []Account, err error) {
	// Only admins can view other accounts.
	if auth.Role != "admin" {
		// When a user is not a "admin", return a list with the same account.
		err = s.db.Model(&accs).Where("id = ?", auth.ID).Offset(offset).Limit(pageSize).Find(&accs).Error
		return
	}

	// Role admin.
	// Get all accounts. Normally.
	err = s.db.Model(&accs).Offset(offset).Limit(pageSize).Find(&accs).Error
	return
}

// UpdateOrCreateAccount Update or create a new account.
func (s Service) UpdateOrCreateAccount(auth Account, id string, payload Account) (err error) {
	// Database account and cache key.
	var acc Account
	var cacheKey = fmt.Sprintf("account_id:%s", id)

	// Only admins can edit others accounts.
	// if the auth account is not a "admin" and want to change other
	if auth.Role != "admin" && auth.ID != id {
		return ErrOnlyAdmin
	}

	// Try to parse ID as uuid.
	if _, err := uuid.Parse(id); err != nil {
		return ErrInconsistentIDs
	}

	// Check if account exists.
	result := s.db.Model(&acc).Where("id = ?", id).First(&acc)

	if result.RowsAffected == 0 {
		// Check needed fields in payload.
		if payload.Name == "" || payload.Email == "" || payload.Password == "" {
			return ErrFieldsRequired
		}
		// If not exists, create.
		if err = s.db.Model(&payload).Create(&payload).Error; err != nil {
			return
		}
	}

	// Update account in database.
	if err = s.db.Model(&payload).Where("id=?", id).Save(&payload).Error; err != nil {
		return
	}

	// Set updated account in cache memory.
	s.cache.Set(cacheKey, payload, cache.DefaultExpiration)
	return
}

// UpdateAccount update specific account fields.
func (s Service) UpdateAccount(auth Account, id string, payload Account) (err error) {
	// Database account and cache key.
	var acc Account
	var cacheKey = fmt.Sprintf("account_id:%s", id)

	// Only admins can edit others accounts.
	// if the auth account is not a "admin" and want to change other
	if auth.Role != "admin" && auth.ID != id {
		return ErrOnlyAdmin
	}

	// Try to parse ID as uuid.
	if _, err := uuid.Parse(id); err != nil {
		return ErrInconsistentIDs
	}

	// Check if account exists.
	err = s.db.Model(&acc).First(&acc, id).Error
	if err != nil {
		return ErrNotFound
	}

	// Update account with fields is not blank/null.
	if err = s.db.Model(&payload).Updates(&payload).Error; err != nil {
		return err
	}

	// Send account to memory cache.
	s.cache.Set(cacheKey, payload, cache.DefaultExpiration)
	return nil
}

// DeleteAccount delete specific account by ID.
func (s Service) DeleteAccount(auth Account, id string) (err error) {
	// Database account and cache key.
	var acc Account
	var cacheKey = fmt.Sprintf("account_id:%s", id)

	// Only admin can delete another accounts.
	if auth.Role != "admin" {
		return ErrOnlyAdmin
	}

	// Admin has the power to delete, but not yourself.
	if auth.Role == "admin" && auth.ID == id {
		return ErrCanNotDeleteYourSelf
	}

	// Delete account in database.
	if err = s.db.Model(&acc).Where("id = ?", id).Delete(&acc).Error; err != nil {
		return err
	}

	// Delete account in cache.
	s.cache.Delete(cacheKey)
	return nil
}

/*
	URL service functions.
*/

// AddURL create a new short URL.
func (s Service) AddURL(auth Account, payload URL) (url URL, err error) {
	// Check if necessary fields was sended.
	if payload.Keyword == "" || payload.URL == "" || payload.Title == "" {
		return payload, errors.New("fields required: keyword, url and title")
	}

	// Check if keyword exists.
	result := s.db.Model(&url).Limit(1).Where("keyword=?", payload.Keyword).Find(&url)
	if result.RowsAffected > 0 {
		return payload, ErrAlreadyExists
	}

	payload.ID = uuid.New().String()
	payload.AccountID = auth.ID
	url = payload

	// Create a new
	if err = s.db.Model(&url).Create(&url).Error; err != nil {
		return
	}

	// Store new URL in memory cache.
	cacheKey := fmt.Sprintf("url_id:%s", url.ID)
	s.cache.Set(cacheKey, url, cache.DefaultExpiration)
	return url, nil
}

// FindURLByID search a specific URL by ID.
func (s Service) FindURLByID(auth Account, id string) (url URL, err error) {
	// Only admin can get any URL.
	if auth.Role == "admin" {
		err = s.db.Model(&url).Where("id = ?", id).First(&url).Error
	} else {
		err = s.db.Model(&url).Where("id = ?", id).Where("account_id=?", auth.ID).First(&url).Error
	}
	// TODO: get item from cache.
	return
}

// FindURLs get a URL list from database.
func (s Service) FindURLs(auth Account, offset, pageSize int) (urls []URL, err error) {
	// Only admin can get all URLs.
	if auth.Role == "admin" {
		err = s.db.Model(&urls).Offset(offset).Limit(pageSize).Find(&urls).Error
	} else {
		err = s.db.Model(&urls).Where("account_id=?", auth.ID).Offset(offset).Limit(pageSize).Find(&urls).Error
	}
	return
}

// UpdateOrCreateURL update or create a url.
func (s Service) UpdateOrCreateURL(auth Account, id string, payload URL) (err error) {
	var cacheKey = fmt.Sprintf("url_id:%s", id)

	// Check if id is a valid UUID.
	_, err = uuid.Parse(id)
	if err != nil {
		return ErrInconsistentIDs
	}

	/*
		TODO: receive URL with account object.
	*/

	// Only admin can edit any URL.
	if auth.Role == "admin" {
		err = s.db.Model(&payload).Where("id = ?", id).Updates(&payload).Error
		if err == gorm.ErrRecordNotFound {
			payload.AccountID = auth.ID
			if err = checkURLFields(payload); err != nil {
				return err
			}
			err = s.db.Model(&payload).Create(&payload).Error
		}
		// Send object to memory cache.
		s.cache.Set(cacheKey, payload.ID, cache.DefaultExpiration)
		return
	}

	// Normal user
	err = s.db.Model(&payload).Where("id = ?", id).Where("account_id = ?", auth.ID).Updates(&payload).Error
	if err == gorm.ErrRecordNotFound {
		payload.AccountID = auth.ID
		if err = checkURLFields(payload); err != nil {
			return err
		}
		err = s.db.Model(&payload).Create(&payload).Error
	}
	// Send object to memory cache.
	s.cache.Set(cacheKey, payload.ID, cache.DefaultExpiration)
	return
}

// UpdateURL update specific URL by ID.
func (s Service) UpdateURL(account Account, id string, reqURL URL) error {
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

// DeleteURL delete account by ID.
func (s Service) DeleteURL(account Account, id string) error {
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

func generateJWT(secretKey string, a Account) (hash string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["id"] = a.ID
	claims["email"] = a.Email
	claims["role"] = a.Role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	hash, err = token.SignedString([]byte(secretKey))
	if err != nil {
		err = errors.New("error to generate JWT: " + err.Error())
		return
	}

	return
}

func checkURLFields(url URL) (err error) {
	if url.Keyword == "" || url.Title == "" || url.URL == "" {
		err = errors.New("fields required: keyword, title and url")
	}
	return
}
