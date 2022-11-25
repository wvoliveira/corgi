package password

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/eko/gocache/v3/cache"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Service encapsulates the authentication logic.
type Service interface {
	Login(ctx context.Context, identity entity.Identity) (model.Session, error)
	Register(ctx context.Context, identity entity.Identity) error

	NewHTTP(r *mux.Router)
	HTTPLogin(w http.ResponseWriter, r *http.Request)
	HTTPRegister(w http.ResponseWriter, r *http.Request)
}

type service struct {
	db    *gorm.DB
	cache *cache.Cache[[]byte]
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB, cache *cache.Cache[[]byte]) Service {
	return service{db, cache}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, identity entity.Identity) (cacheValue model.Session, err error) {
	l := logger.Logger(ctx)

	identityDB := entity.Identity{}
	err = s.db.Model(&entity.Identity{}).Where("provider = ? AND uid = ?", identity.Provider, identity.UID).First(&identityDB).Error
	if err != nil {
		l.Warn().Caller().Msg(err.Error())
		return cacheValue, e.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(identityDB.Password), []byte(identity.Password)); err != nil {
		l.Info().Caller().Msg("authentication failed")
		return cacheValue, e.ErrUnauthorized
	}

	// Get user info from database or return an error if not exists.
	user := entity.User{}
	err = s.db.Model(&entity.User{}).Where("id = ?", identityDB.UserID).First(&user).Error
	if err != nil {
		l.Error().Caller().Msg(err.Error())
		return cacheValue, err
	}

	tokenAccess := uuid.New().String()
	tokenRefresh := uuid.New().String()

	// Store tokens and user struct in cache.
	// And expires this session in 5 minutes.
	createdAt := time.Now()
	expiresIn := createdAt.Add(time.Minute * 5)

	cacheKey := fmt.Sprintf("token_%s", tokenAccess)
	cacheValue = model.Session{
		ID:           cacheKey,
		CreatedAt:    createdAt,
		TokenAccess:  tokenAccess,
		TokenRefresh: tokenRefresh,
		ExpiresIn:    expiresIn,
		User:         user,
	}
	err = s.cache.Set(ctx, cacheKey, cacheValue.Encode())
	return
}

// Register a new user to our database.
func (s service) Register(ctx context.Context, identity entity.Identity) (err error) {
	l := logger.Logger(ctx)

	user := entity.User{}
	identityDB := entity.Identity{}

	err = s.db.Model(&entity.Identity{}).
		Where("provider = ? AND uid = ?", identity.Provider, identity.UID).
		First(&identityDB).Error

	if err == nil {
		l.Warn().Caller().Msg(fmt.Sprintf("provider '%s' and uid '%s' already exists", identity.Provider, identity.UID))
		return e.ErrAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(identity.Password), 8)
	if err != nil {
		l.Error().Caller().Msg(err.Error())
		return e.ErrInternalServerError
	}

	identity.ID = uuid.New().String()
	identity.CreatedAt = time.Now()
	identity.Password = string(hashedPassword)

	t := true
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.Role = "user"
	user.Active = &t
	user.Identities = append(user.Identities, identity)

	err = s.db.Model(&entity.User{}).Create(&user).Error
	if err != nil {
		l.Error().Caller().Msg(err.Error())
	}
	return
}
