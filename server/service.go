package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/mail"
	"time"

	"github.com/go-kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service store all methods. Yeah monolithic.
type Service struct {
	logger log.Logger
	ctx    context.Context
	db     *redis.Client
	cache  *redis.Client
	secret string
}

// NewService create a new service with database and cache.
func NewService(logger log.Logger, ctx context.Context, secretKey string, db *redis.Client, cache *redis.Client) Service {
	return Service{
		logger: logger,
		ctx:    ctx,
		db:     db,
		cache:  cache,
		secret: secretKey,
	}
}

// SignIn login with email and password.
func (s Service) SignIn(payload Account) (account Account, err error) {
	var (
		found, foundInCache bool
		keys                []string
	)

	dbKeyPattern, cacheKeyPattern := s.getAccountKey("*", payload.Email)

	_, err = mail.ParseAddress(payload.Email)
	if err != nil {
		return account, ErrEmailNotValid
	}

	if payload.Password == "" {
		return account, ErrFieldsRequired
	}

	keys, err = s.cache.Keys(s.ctx, cacheKeyPattern).Result()
	if len(keys) != 0 {
		foundInCache = true
	}

	if foundInCache {
		item, err := s.cache.Get(s.ctx, keys[0]).Result()

		if err == nil {
			if err = json.Unmarshal([]byte(item), &account); err != nil {
				return account, err
			}
		}
	}

	if !foundInCache {
		keys, err = s.db.Keys(s.ctx, dbKeyPattern).Result()
		if len(keys) != 0 {
			found = true
		}
	}

	if found {
		item, err := s.db.Get(s.ctx, keys[0]).Result()

		// Not found in database.
		if err == redis.Nil {
			return account, ErrUnauthorized
		}

		if err == nil {
			if err = json.Unmarshal([]byte(item), &account); err != nil {
				return account, err
			}
		}
	}

	if err != redis.Nil && err != nil {
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(payload.Password)); err != nil {
		return account, ErrUnauthorized
	}

	tokenHash, err := generateJWT(s.secret, account)
	if err != nil {
		return
	}

	account.Token = tokenHash
	accountJs, err := json.Marshal(account)
	if err != nil {
		return
	}

	if !foundInCache {
		_, cacheKey := s.getAccountKey(account.ID, payload.Email)
		_ = s.cache.Set(s.ctx, cacheKey, accountJs, 0).Err()
	}
	return
}

// SignUp register with e-mail and password.
func (s Service) SignUp(payload Account) (err error) {
	var (
		dbKey, cacheKey string
		accountJs       []byte
		keys            []string
	)

	dbKeyPattern, cacheKeyPattern := s.getAccountKey("*", payload.Email)

	_, err = mail.ParseAddress(payload.Email)
	if err != nil {
		return ErrEmailNotValid
	}

	if payload.Password == "" {
		return ErrFieldsRequired
	}

	keys, _ = s.cache.Keys(s.ctx, cacheKeyPattern).Result()
	if len(keys) > 0 {
		return ErrAlreadyExists
	}

	keys, err = s.db.Keys(s.ctx, dbKeyPattern).Result()
	if err != redis.Nil && err != nil {
		return
	}
	if len(keys) > 0 {
		return ErrAlreadyExists
	}

	// Salt and hash the password using the bcrypt algorithm.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 8)
	if err != nil {
		return ErrInternalServerError
	}

	// Create a random ID and default role for new user.
	payload.ID = uuid.New().String()
	payload.CreatedAt = time.Now()
	payload.UpdatedAt = time.Now()

	payload.Password = string(hashedPassword)
	payload.Role = "user"
	payload.Active = "true"

	accountJs, err = json.Marshal(payload)
	if err != nil {
		return err
	}

	dbKey, cacheKey = s.getAccountKey(payload.ID, payload.Email)
	if err = s.db.Set(s.ctx, dbKey, accountJs, 0).Err(); err != nil {
		return err
	}
	_ = s.cache.Set(s.ctx, cacheKey, accountJs, 10*time.Minute).Err()
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
