package server

import (
	"encoding/json"
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
		found     bool
		accountJs []byte
	)

	dbKey, cacheKey := s.getAccountKey("*", payload.Email)

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

	_, err = s.cache.Get(s.ctx, cacheKey).Result()
	if err == nil {
		return account, ErrAlreadyExists
	}

	if err == redis.Nil {
		_, err = s.db.Get(s.ctx, dbKey).Result()

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

	dbKey, cacheKey = s.getAccountKey(payload.ID, payload.Email)

	if err = s.db.Set(s.ctx, dbKey, accountJs, 0).Err(); err != nil {
		return
	}

	_ = s.cache.Set(s.ctx, cacheKey, accountJs, 10).Err()

	account = payload
	return
}

// FindAccountByID find a account with specific ID.
func (s Service) FindAccountByID(auth Account, id string) (account Account, err error) {
	var (
		dbKey, cacheKey     string
		found, foundInCache bool
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
		foundInCache = true
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
		return account, ErrNotFound
	}

	if err = json.Unmarshal([]byte(item), &account); err != nil {
		return
	}

	if !foundInCache {
		_, cacheKey = s.getAccountKey(id, account.Email)
		_ = s.cache.Set(s.ctx, cacheKey, item, 10).Err()
	}
	return
}

// FindAccounts Get a list of accounts.
func (s Service) FindAccounts(auth Account, offset, pageSize int) (accounts []Account, err error) {
	var (
		dbKey   string
		account Account
	)

	if auth.Role == "admin" {
		dbKey, _ = s.getAccountKey("*", "*")
	} else {
		dbKey, _ = s.getAccountKey(auth.ID, auth.Email)
	}

	items, err := s.db.Sort(s.ctx, dbKey, &redis.Sort{Offset: int64(offset), Count: int64(pageSize), Order: "ASC"}).Result()
	if err != nil {
		return
	}

	for _, item := range items {
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
		dbKey, cacheKey = s.getAccountKey(id, "*")
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
		payload.ID = id
		payload.CreatedAt = time.Now()
		payload.UpdatedAt = time.Now()
		payload.Active = "true"

		accountJs, err = json.Marshal(payload)
		if err != nil {
			return
		}

		dbKey, cacheKey = s.getAccountKey(id, payload.Email)

		if err = s.db.Set(s.ctx, dbKey, accountJs, 0).Err(); err != nil {
			return
		}

		_ = s.cache.Set(s.ctx, cacheKey, accountJs, 10).Err()
		return
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

	dbKey, cacheKey = s.getAccountKey(id, payload.Email)

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

	item, err := s.cache.Get(s.ctx, cacheIDKey).Result()
	if err == nil {
		found = true
	}

	// If not found in cache memory, search in database (more slowly).
	if !found {
		item, err = s.db.Get(s.ctx, dbIDKey).Result()

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

	dbIDKey, dbEmailKey := s.dbAccountKey(id, account.Email)
	cacheIDKey, cacheEmailKey := s.cacheAccountKey(id, account.Email)

	// Send a new account into database.
	if err = sendPayloadPipelineDB(s.db, s.ctx, dbIDKey, dbEmailKey, accountJs); err != nil {
		return
	}

	// Send a new account to cache. No problem if error apper.
	_ = sendPayloadPipelineCache(s.cache, s.ctx, cacheIDKey, cacheEmailKey, accountJs)
	return
}
