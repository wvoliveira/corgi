package server

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

/*
	I'm using redis for store account and URL. So, theses are patterns for keys:
	Database: db_url_id:<url id>:url_keyword:<url keyword>:account_id:<account id>
	Cache: cache_url_id:<url id>:url_keyword:<url keyword>:account_id:<account id>
*/

func (s Service) getURLKey(id, keyword, accountID string) (dbKey, cacheKey string) {
	return fmt.Sprintf("db_url_id:%s:url_keyword:%s:account_id:%s", id, keyword, accountID),
		fmt.Sprintf("cache_url_id:%s:url_keyword:%s:account_id:%s", id, keyword, accountID)
}

// AddURL create a new short URL.
func (s Service) AddURL(auth Account, payload URL) (url URL, err error) {
	var (
		found bool
		urlJs []byte
	)

	dbKey, cacheKey := s.getURLKey("*", payload.Keyword, "*")

	if err = checkURLFields(payload); err != nil {
		return
	}

	// Check if URL exist in cache.
	_, err = s.cache.Get(s.ctx, cacheKey).Result()
	if err == nil {
		return url, ErrAlreadyExists
	}

	// If not found in cache memory, search in database (more slowly).
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
		return url, ErrAlreadyExists
	}

	payload.ID = uuid.New().String()
	payload.CreatedAt = time.Now()
	payload.UpdatedAt = time.Now()
	payload.AccountID = auth.ID
	payload.Active = "true"

	urlJs, err = json.Marshal(payload)
	if err != nil {
		return
	}

	dbKey, cacheKey = s.getURLKey(payload.ID, payload.Keyword, auth.ID)

	if err = s.db.Set(s.ctx, dbKey, urlJs, 0).Err(); err != nil {
		return
	}

	_ = s.cache.Set(s.ctx, cacheKey, urlJs, 10).Err()

	url = payload
	return
}

// FindURLByID search a specific URL by ID.
func (s Service) FindURLByID(auth Account, id string) (url URL, err error) {
	var (
		dbKey, cacheKey     string
		found, foundInCache bool
	)

	_, err = uuid.Parse(id)
	if err != nil {
		return url, ErrInconsistentIDs
	}

	if auth.Role == "admin" {
		dbKey, cacheKey = s.getURLKey(id, "*", "*")
	} else {
		dbKey, cacheKey = s.getURLKey(id, "*", auth.ID)
	}

	item, err := s.cache.Get(s.ctx, cacheKey).Result()
	if err == nil {
		found = true
		foundInCache = true
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
		return url, ErrNotFound
	}

	if err = json.Unmarshal([]byte(item), &url); err != nil {
		return
	}

	// Send a new account to cache. No problem if error apper.
	if !foundInCache {
		_, cacheKey = s.getURLKey(id, url.Keyword, auth.ID)
		_ = s.cache.Set(s.ctx, cacheKey, item, 10).Err()
	}
	return
}

// FindURLs get a URL list from database.
func (s Service) FindURLs(auth Account, offset, pageSize int) (urls []URL, err error) {
	var (
		dbKey string
		url   URL
	)

	if auth.Role == "admin" {
		dbKey, _ = s.getURLKey("*", "*", "*")
	} else {
		dbKey, _ = s.getURLKey("*", "*", auth.ID)
	}

	items, err := s.db.Sort(s.ctx, dbKey, &redis.Sort{Offset: int64(offset), Count: int64(pageSize), Order: "ASC"}).Result()
	if err != nil {
		return
	}

	for _, item := range items {
		err = json.Unmarshal([]byte(item), &url)
		if err != nil {
			return
		}
		urls = append(urls, url)
	}
	return
}

// UpdateOrCreateURL update or create a url.
func (s Service) UpdateOrCreateURL(auth Account, id string, payload URL) (err error) {
	var (
		dbKey, cacheKey string
		found           bool
		url             URL
		urlJs           []byte
	)

	_, err = uuid.Parse(id)
	if err != nil {
		return ErrInconsistentIDs
	}

	if auth.Role == "admin" {
		dbKey, cacheKey = s.getURLKey("*", "*", "*")
	} else {
		dbKey, cacheKey = s.getURLKey("*", "*", auth.ID)
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
		payload.AccountID = auth.ID
		payload.Active = "true"

		urlJs, err = json.Marshal(payload)
		if err != nil {
			return
		}

		dbKey, cacheKey = s.getURLKey(id, payload.Keyword, auth.ID)

		if err = s.db.Set(s.ctx, dbKey, urlJs, 0).Err(); err != nil {
			return
		}

		_ = s.cache.Set(s.ctx, cacheKey, urlJs, 10).Err()
		return
	}

	err = json.Unmarshal([]byte(item), &url)
	if err != nil {
		return
	}

	url.URL = payload.URL
	url.Title = payload.Title
	url.UpdatedAt = time.Now()

	urlJs, err = json.Marshal(url)
	if err != nil {
		return
	}

	dbKey, cacheKey = s.getURLKey(id, payload.Keyword, url.AccountID)

	if err = s.db.Set(s.ctx, dbKey, urlJs, 0).Err(); err != nil {
		return
	}

	_ = s.cache.Set(s.ctx, cacheKey, urlJs, 10).Err()
	return
}

// UpdateURL update specific URL by ID.
func (s Service) UpdateURL(auth Account, id string, payload URL) (err error) {
	var (
		dbKey, cacheKey string
		found           bool
		url             URL
		urlJs           []byte
	)

	_, err = uuid.Parse(id)
	if err != nil {
		return ErrInconsistentIDs
	}

	if auth.Role == "admin" {
		dbKey, cacheKey = s.getURLKey(id, "*", "*")
	} else {
		dbKey, cacheKey = s.getURLKey(id, "*", auth.ID)
	}

	// Check if account exist in cache.
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

	err = json.Unmarshal([]byte(item), &url)
	if err != nil {
		return
	}

	url.URL = payload.URL
	url.Title = payload.Title
	url.UpdatedAt = time.Now()

	urlJs, err = json.Marshal(url)
	if err != nil {
		return
	}

	dbKey, cacheKey = s.getURLKey(id, url.Keyword, url.AccountID)

	if err = s.db.Set(s.ctx, dbKey, urlJs, 0).Err(); err != nil {
		return
	}

	_ = s.cache.Set(s.ctx, cacheKey, urlJs, 10).Err()
	return
}

// DeleteURL delete URL by ID.
func (s Service) DeleteURL(auth Account, id string) (err error) {
	var (
		dbKey, cacheKey string
		found           bool
		url             URL
		urlJs           []byte
	)

	_, err = uuid.Parse(id)
	if err != nil {
		return ErrInconsistentIDs
	}

	if auth.Role == "admin" {
		dbKey, cacheKey = s.getURLKey(id, "*", "*")
	} else {
		dbKey, cacheKey = s.getURLKey(id, "*", auth.ID)
	}

	// Check if account exist in cache.
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

	err = json.Unmarshal([]byte(item), &url)
	if err != nil {
		return
	}

	url.Active = "false"
	url.UpdatedAt = time.Now()

	urlJs, err = json.Marshal(url)
	if err != nil {
		return
	}

	dbKey, cacheKey = s.getURLKey(id, url.Keyword, url.AccountID)

	if err = s.db.Set(s.ctx, dbKey, urlJs, 0).Err(); err != nil {
		return
	}

	_ = s.cache.Set(s.ctx, cacheKey, urlJs, 10).Err()
	return
}
