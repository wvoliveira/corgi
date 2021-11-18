package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/mail"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service store all methods. Yeah monolithic.
type Service struct {
	ctx    context.Context
	db     *redis.Client
	cache  *redis.Client
	secret string
}

// NewService create a new service with database and cache.
func NewService(ctx context.Context, secretKey string, db *redis.Client, cache *redis.Client) Service {
	return Service{
		ctx:    ctx,
		db:     db,
		cache:  cache,
		secret: secretKey,
	}
}

// SignIn login with email and password.
func (s Service) SignIn(payload Account) (account Account, err error) {
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
	item, err := s.cache.Get(s.ctx, cacheEmailKey).Result()
	if err == nil {
		if err = json.Unmarshal([]byte(item), &account); err != nil {
			return
		}
	}

	// If not found in cache memory, search in database (more slowly).
	if err == redis.Nil {
		item, err = s.db.Get(s.ctx, dbEmailKey).Result()

		// Not found in database.
		if err == redis.Nil {
			return account, ErrUnauthorized
		}

		// Found! Try to unmarshal and set to account variable.
		if err == nil {
			if err = json.Unmarshal([]byte(item), &account); err != nil {
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
	_ = s.cache.Set(s.ctx, cacheEmailKey, accountJs, 0).Err()
	return
}

// SignUp register with e-mail and password.
func (s Service) SignUp(payload Account) (err error) {
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
	_, err = s.cache.Get(s.ctx, cacheEmailKey).Result()
	if err == nil {
		return ErrAlreadyExists

		// Get from database here.
	} else {
		_, err = s.db.Get(s.ctx, dbEmailKey).Result()
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
	if _, err = s.db.Pipelined(s.ctx, func(rdb redis.Pipeliner) error {
		rdb.Set(s.ctx, dbIDKey, accountJs, 0)
		rdb.Set(s.ctx, dbEmailKey, accountJs, 0)
		return nil
	}); err != nil {
		return
	}

	// Send a new account to cache. No problem if error apper.
	_, _ = s.cache.Pipelined(s.ctx, func(rdb redis.Pipeliner) error {
		rdb.Set(s.ctx, cacheIDKey, accountJs, 10)
		rdb.Set(s.ctx, cacheEmailKey, accountJs, 10)
		return nil
	})
	return
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

func checkAccountRequest(auth Account, id string, payload Account) (err error) {
	if payload.Email != "" {
		if _, err = mail.ParseAddress(payload.Email); err != nil {
			return ErrEmailNotValid
		}
	}

	if auth.Role != "admin" && auth.ID != id {
		return ErrOnlyAdmin
	}

	if _, err := uuid.Parse(id); err != nil {
		return ErrInconsistentIDs
	}
	return
}

func sendPayloadPipelineDB(db *redis.Client, ctx context.Context, keyID, keyEmail string, value interface{}) (err error) {
	if _, err = db.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Set(ctx, keyID, value, 0)
		rdb.Set(ctx, keyEmail, value, 0)
		return nil
	}); err != nil {
		return
	}
	return
}

func sendPayloadPipelineCache(cache *redis.Client, ctx context.Context, keyID, keyEmail string, value interface{}) (err error) {
	// Send a new account to cache. No problem if error apper.
	_, err = cache.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.Set(ctx, keyID, value, 10)
		rdb.Set(ctx, keyEmail, value, 10)
		return nil
	})
	return
}
