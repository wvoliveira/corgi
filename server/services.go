package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Service store all methods. Yeah monolithic.
type Service struct {
	db     *redis.Client
	cache  *redis.Client
	secret string
}

// NewService create a new service with database and cache.
func NewService(secretKey string, db *redis.Client, cache *redis.Client) Service {
	return Service{
		db:     db,
		cache:  cache,
		secret: secretKey,
	}
}

func (s Service) dbAccountKey(id, email string) (_id, _email string) {
	return fmt.Sprintf("db_account_id:%s", id), fmt.Sprintf("db_account_email:%s", email)
}

func (s Service) cacheAccountKey(id, email string) (_id, _email string) {
	return fmt.Sprintf("cache_account_id:%s", id), fmt.Sprintf("cache_account_email:%s", email)
}

// SignIn login with email and password.
func (s Service) SignIn(ctx context.Context, payload Account) (account Account, err error) {
	_, dbEmailKey := s.dbAccountKey(payload.ID, payload.Email)
	_, cacheEmailKey := s.cacheAccountKey(payload.ID, payload.Email)

	// Check e-mail pattern.
	_, err = mail.ParseAddress(payload.Email)
	if err != nil {
		return account, ErrEmailNotValid
	}

	// Check if e-mail and password was sended.
	if payload.Password == "" {
		return account, ErrFieldsRequired
	}

	// Check if account exist in cache.
	item, err := s.cache.Get(ctx, cacheEmailKey).Result()
	if err == nil {
		if err = json.Unmarshal([]byte(item), account); err != nil {
			return
		}
	}

	// If not found in cache memory, search in database (more slowly).
	if err == redis.Nil {
		item, err = s.db.Get(ctx, dbEmailKey).Result()

		// Not found in database.
		if err == redis.Nil {
			return account, ErrUnauthorized
		}

		// Found! Try to unmarshal and set to account variable.
		if err == nil {
			if err = json.Unmarshal([]byte(item), account); err != nil {
				return
			}
		}
	}

	// If err is unknown, we are in trouble.
	if err != redis.Nil && err != nil {
		return
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(payload.Password)); err != nil {
		return account, ErrUnauthorized
	}

	// Generate JWT with specific claims.
	tokenHash, err := generateJWT(s.secret, account)
	if err != nil {
		return
	}

	// Set token account with JWT.
	account.Token = tokenHash

	accountJs, err := json.Marshal(account)
	if err != nil {
		return
	}

	// Try to store the account in cache and not problem if cache is unavailable.
	_ = s.cache.Set(ctx, cacheEmailKey, accountJs, 0).Err()
	return
}

// SignUp register with e-mail and password.
func (s Service) SignUp(ctx context.Context, payload Account) (err error) {
	dbIDKey, dbEmailKey := s.dbAccountKey(payload.ID, payload.Email)
	cacheIDKey, cacheEmailKey := s.cacheAccountKey(payload.ID, payload.Email)

	// Check e-mail pattern.
	_, err = mail.ParseAddress(payload.Email)
	if err != nil {
		return ErrEmailNotValid
	}

	if payload.Password == "" {
		return ErrFieldsRequired
	}

	// Check if account exist in cache.
	_, err = s.cache.Get(ctx, cacheEmailKey).Result()
	if err == nil {
		return ErrAlreadyExists

		// Get from database here.
	} else {
		_, err = s.db.Get(ctx, dbEmailKey).Result()
		if err == nil {
			return ErrAlreadyExists
		}

		// We're in trouble.
		if err != redis.Nil && err != nil {
			return
		}
	}

	// Salt and hash the password using the bcrypt algorithm.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 8)
	if err != nil {
		return ErrInternalServerError
	}

	// Create a random ID and default role for new user.
	payload.ID = uuid.New().String()
	payload.Role = "user"
	payload.Password = string(hashedPassword)

	accountJs, err := json.Marshal(payload)
	if err != nil {
		return
	}

	// Send a new account into database.
	if _, err = s.db.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Set(ctx, dbIDKey, accountJs, 0)
		rdb.Set(ctx, dbEmailKey, accountJs, 0)
		return nil
	}); err != nil {
		return
	}

	// Send a new account to cache. No problem if error apper.
	_, _ = s.db.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Set(ctx, cacheIDKey, accountJs, 10)
		rdb.Set(ctx, cacheEmailKey, accountJs, 10)
		return nil
	})
	return
}

// AddAccount create a new account.
func (s Service) AddAccount(ctx context.Context, auth, payload Account) (account Account, err error) {
	dbIDKey, dbEmailKey := s.dbAccountKey(payload.ID, payload.Email)
	cacheIDKey, cacheEmailKey := s.cacheAccountKey(payload.ID, payload.Email)

	// Check e-mail pattern.
	_, err = mail.ParseAddress(payload.Email)
	if err != nil {
		return account, ErrEmailNotValid
	}

	// Only admin can create a new account without signup process.
	// TODO: create more roles without hardcoded.
	if auth.Role != "admin" {
		return account, ErrOnlyAdmin
	}

	if payload.Password == "" {
		return account, ErrFieldsRequired
	}

	// Check if account exist in cache.
	_, err = s.cache.Get(ctx, cacheEmailKey).Result()
	if err == nil {
		return account, ErrAlreadyExists
	}

	var found bool

	// If not found in cache memory, search in database (more slowly).
	if err == redis.Nil {
		_, err = s.db.Get(ctx, dbEmailKey).Result()

		// Not found in database.
		if err == redis.Nil {
			found = false
		} else if err == nil {
			found = true
		} else {
			return
		}
	}

	if found {
		return account, ErrAlreadyExists
	}

	// Set a random ID and create account.
	payload.ID = uuid.New().String()

	accountJs, err := json.Marshal(payload)
	if err != nil {
		return
	}

	// Send a new account into database.
	if _, err = s.db.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Set(ctx, dbIDKey, accountJs, 0)
		rdb.Set(ctx, dbEmailKey, accountJs, 0)
		return nil
	}); err != nil {
		return
	}

	// Send a new account to cache. No problem if error apper.
	_, _ = s.db.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Set(ctx, cacheIDKey, accountJs, 10)
		rdb.Set(ctx, cacheEmailKey, accountJs, 10)
		return nil
	})

	account = payload
	return
}

// FindAccountByID find a account with specific ID.
func (s Service) FindAccountByID(ctx context.Context, auth Account, id string) (account Account, err error) {
	dbIDKey, _ := s.dbAccountKey(id, "")
	cacheIDKey, _ := s.cacheAccountKey(id, "")

	// Only admin can view another accounts.
	// TODO: create more roles without hardcoded.
	if auth.Role != "admin" && auth.ID != id {
		return account, ErrOnlyAdmin
	}

	var found, foundInCache bool

	// Check if account exist in cache.
	item, err := s.cache.Get(ctx, cacheIDKey).Result()
	if err == nil {
		found = true
		foundInCache = true
	}

	// If not found in cache memory, search in database (more slowly).
	if !found {
		item, err = s.db.Get(ctx, dbIDKey).Result()

		// Not found in database.
		if err == redis.Nil {
			found = false
		} else if err == nil {
			found = true
		} else {
			return
		}
	}

	if !found {
		return account, ErrNotFound
	}

	if err = json.Unmarshal([]byte(item), account); err != nil {
		return
	}

	// Send a new account to cache. No problem if error apper.
	if !foundInCache {
		cacheIDKey, cacheEmailKey := s.cacheAccountKey(id, account.Email)

		_, _ = s.db.Pipelined(ctx, func(rdb redis.Pipeliner) error {
			rdb.Set(ctx, cacheIDKey, item, 10)
			rdb.Set(ctx, cacheEmailKey, item, 10)
			return nil
		})
	}
	return
}

// FindAccounts Get a list of accounts.
func (s Service) FindAccounts(ctx context.Context, auth Account, offset, pageSize int) (accounts []Account, err error) {
	dbIDKey, _ := s.dbAccountKey("*", "")
	dbAuthIDKey, _ := s.dbAccountKey(auth.ID, "")

	var items []string

	// TODO: Cache response.
	if auth.Role != "admin" {
		// When a user is not a "admin", return a list with the same account.
		items, err = s.db.Sort(ctx, dbAuthIDKey, &redis.Sort{Offset: int64(offset), Count: int64(pageSize), Order: "ASC"}).Result()
		if err != nil {
			return
		}
	} else {
		items, err = s.db.Sort(ctx, dbIDKey, &redis.Sort{Offset: int64(offset), Count: int64(pageSize), Order: "ASC"}).Result()
		if err != nil {
			return
		}
	}

	for _, item := range items {
		var account Account
		err = json.Unmarshal([]byte(item), account)
		if err != nil {
			return
		}
		accounts = append(accounts, account)
	}

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
