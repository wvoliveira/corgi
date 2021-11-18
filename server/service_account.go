package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) getAccountKey(id, email string) (dbKey, cacheKey string) {
	return fmt.Sprintf("db_account_id:%s:account_email:%s", id, email),
		fmt.Sprintf("cache_account_id:%s:account_email:%s", id, email)
}

// AddAccount create a new account.
func (s Service) AddAccount(auth, payload Account) (account Account, err error) {
	var (
		accountJs           []byte
		keys                []string
		found, foundInCache bool
	)

	dbKeyPattern, cacheKeyPattern := s.getAccountKey("*", payload.Email)

	_, err = mail.ParseAddress(payload.Email)
	if err != nil {
		return account, ErrEmailNotValid
	}

	if auth.Role != "admin" {
		return account, ErrOnlyAdmin
	}

	if payload.Password == "" {
		return account, ErrFieldsRequired
	}

	keys, _ = s.cache.Keys(s.ctx, cacheKeyPattern).Result()
	if len(keys) > 0 {
		return account, ErrAlreadyExists
	}

	keys, err = s.db.Keys(s.ctx, dbKeyPattern).Result()
	if err != redis.Nil && err != nil {
		return account, err
	}
	if len(keys) > 0 {
		found = true
	}

	if found && !foundInCache {
		item, err := s.db.Get(s.ctx, keys[0]).Result()
		if err != redis.Nil && err != nil {
			return account, ErrAlreadyExists
		}

		_ = s.cache.Set(s.ctx, keys[0], item, 10*time.Minute).Err()
		return account, ErrAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 8)
	if err != nil {
		return account, ErrInternalServerError
	}

	payload.ID = uuid.New().String()
	payload.Password = string(hashedPassword)
	payload.CreatedAt = time.Now()
	payload.UpdatedAt = time.Now()
	payload.Active = "true"

	accountJs, err = json.Marshal(payload)
	if err != nil {
		return
	}

	dbKey, cacheKey := s.getAccountKey(payload.ID, payload.Email)

	if err = s.db.Set(s.ctx, dbKey, accountJs, 0).Err(); err != nil {
		return
	}
	_ = s.cache.Set(s.ctx, cacheKey, accountJs, 10).Err()
	return payload, nil
}

// FindAccountByID find a account with specific ID.
func (s Service) FindAccountByID(auth Account, id string) (account Account, err error) {
	var (
		dbKeyPattern, cacheKeyPattern string
		found, foundInCache           bool
		keys                          []string
		item                          string
	)

	if err = checkAccountRequest(auth, id, Account{}); err != nil {
		return
	}

	if auth.Role == "admin" {
		dbKeyPattern, cacheKeyPattern = s.getAccountKey(id, "*")
	} else {
		dbKeyPattern, cacheKeyPattern = s.getAccountKey(id, auth.Email)
	}

	keys, _ = s.cache.Keys(s.ctx, cacheKeyPattern).Result()
	if len(keys) > 0 {
		found = true
		foundInCache = true
	}

	if !found {
		keys, _ = s.db.Keys(s.ctx, dbKeyPattern).Result()
		if len(keys) > 0 {
			found = true
		}
	}

	if !found && !foundInCache {
		return account, ErrNotFound
	}

	if foundInCache {
		item, err = s.cache.Get(s.ctx, keys[0]).Result()
	} else {
		item, err = s.db.Get(s.ctx, keys[0]).Result()
	}

	if err != redis.Nil && err != nil {
		return account, err
	}

	if err = json.Unmarshal([]byte(item), &account); err != nil {
		return
	}

	if !foundInCache {
		_ = s.cache.Set(s.ctx, keys[0], item, 10).Err()
	}
	return
}

// FindAccounts Get a list of accounts.
func (s Service) FindAccounts(auth Account, offset, pageSize int) (accounts []Account, err error) {
	var (
		dbKeyPattern string
		keys         []string

		items []string
		item  string

		account Account
	)

	if auth.Role == "admin" {
		dbKeyPattern, _ = s.getAccountKey("*", "*")
	} else {
		dbKeyPattern, _ = s.getAccountKey(auth.ID, auth.Email)
	}

	keys, _, err = s.db.Scan(s.ctx, uint64(offset), dbKeyPattern, int64(pageSize)).Result()

	for _, key := range keys {
		item, err = s.db.Get(s.ctx, key).Result()
		if err != redis.Nil && err != nil {
			return accounts, err
		}
		items = append(items, item)
	}

	for _, item = range items {
		err = json.Unmarshal([]byte(item), &account)
		if err != nil {
			return
		}
		accounts = append(accounts, account)
	}
	return
}

// UpdateOrCreateAccount Update or create a new account.
func (s Service) UpdateOrCreateAccount(auth Account, id string, payload Account) (err error) {
	var (
		dbKeyPattern, cacheKeyPattern string
		found, foundInCache           bool
		keys                          []string
		account                       Account
		accountJs                     []byte
		item                          string
	)

	if err = checkAccountRequest(auth, id, payload); err != nil {
		return
	}

	if auth.Role == "admin" {
		dbKeyPattern, cacheKeyPattern = s.getAccountKey(id, "*")
	} else {
		dbKeyPattern, cacheKeyPattern = s.getAccountKey(id, auth.Email)
	}

	keys, _ = s.cache.Keys(s.ctx, cacheKeyPattern).Result()
	if len(keys) > 0 {
		found = true
		foundInCache = true
	}

	if !found {
		keys, _ = s.db.Keys(s.ctx, dbKeyPattern).Result()
		if len(keys) > 0 {
			found = true
		}
	}

	if !found && !foundInCache {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 8)
		if err != nil {
			return errors.New("error to generate a hash password")
		}

		payload.ID = id
		payload.CreatedAt = time.Now()
		payload.UpdatedAt = time.Now()
		payload.Password = string(hashedPassword)
		payload.Role = "user"
		payload.Active = "true"
	} else {
		item, err = s.db.Get(s.ctx, keys[0]).Result()
		if err != redis.Nil && err != nil {
			return errors.New("error to get account from database: " + err.Error())
		}

		if err = json.Unmarshal([]byte(item), &account); err != nil {
			return
		}

		payload.ID = account.ID
		payload.CreatedAt = account.CreatedAt
		payload.UpdatedAt = time.Now()
		payload.LastLogin = account.LastLogin

		payload.Password = account.Password
		payload.Role = account.Role
		payload.Active = account.Active
	}

	accountJs, err = json.Marshal(payload)
	if err != nil {
		return
	}

	dbKey, cacheKey := s.getAccountKey(id, payload.Email)

	if err = s.db.Set(s.ctx, dbKey, accountJs, 0).Err(); err != nil {
		return
	}

	_ = s.cache.Set(s.ctx, cacheKey, accountJs, 10).Err()
	return
}

// UpdateAccount update specific account fields.
func (s Service) UpdateAccount(auth Account, id string, payload Account) (err error) {
	var (
		dbKey, cacheKey string
		found           bool
		account         Account
		accountJs       []byte
	)

	if err = checkAccountRequest(auth, id, payload); err != nil {
		return
	}

	if auth.Role == "admin" {
		dbKey, cacheKey = s.getAccountKey(id, "*")
	} else {
		dbKey, cacheKey = s.getAccountKey(id, auth.Email)
	}

	item, err := s.cache.Get(s.ctx, cacheKey).Result()
	if err == nil {
		found = true
	}

	if !found {
		item, err = s.db.Get(s.ctx, dbKey).Result()

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
		return ErrNotFound
	}

	err = json.Unmarshal([]byte(item), &account)
	if err != nil {
		return
	}

	account.Name = payload.Name
	account.Email = payload.Email
	account.UpdatedAt = time.Now()

	accountJs, err = json.Marshal(account)
	if err != nil {
		return
	}

	dbKey, cacheKey = s.getAccountKey(id, account.Email)

	if err = s.db.Set(s.ctx, dbKey, accountJs, 0).Err(); err != nil {
		return
	}

	_ = s.cache.Set(s.ctx, cacheKey, accountJs, 10).Err()
	return
}

// DeleteAccount delete specific account by ID.
func (s Service) DeleteAccount(auth Account, id string) (err error) {
	var (
		dbKey, cacheKey string
		found           bool
		account         Account
		accountJs       []byte
	)

	if err = checkAccountRequest(auth, id, Account{}); err != nil {
		return
	}

	if auth.Role == "admin" {
		dbKey, cacheKey = s.getAccountKey(id, "*")
	} else {
		dbKey, cacheKey = s.getAccountKey(id, auth.Email)
	}

	item, err := s.cache.Get(s.ctx, cacheKey).Result()
	if err == nil {
		found = true
	}

	// If not found in cache memory, search in database (more slowly).
	if !found {
		item, err = s.db.Get(s.ctx, dbKey).Result()

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
		return ErrNotFound
	}

	err = json.Unmarshal([]byte(item), &account)
	if err != nil {
		return
	}

	account.Active = "false"
	account.UpdatedAt = time.Now()

	accountJs, err = json.Marshal(account)
	if err != nil {
		return
	}

	dbKey, cacheKey = s.getAccountKey(id, account.Email)

	if err = s.db.Set(s.ctx, dbKey, accountJs, 0).Err(); err != nil {
		return
	}

	_ = s.cache.Set(s.ctx, cacheKey, accountJs, 10).Err()
	return
}
