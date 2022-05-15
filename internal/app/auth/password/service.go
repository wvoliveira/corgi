package password

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/elga-io/corgi/internal/app/entity"
	e "github.com/elga-io/corgi/internal/pkg/errors"
	"github.com/elga-io/corgi/internal/pkg/jwt"
	"github.com/elga-io/corgi/internal/pkg/logger"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Service encapsulates the authentication logic.
type Service interface {
	Login(ctx context.Context, identity entity.Identity) (entity.Token, entity.Token, error)
	Register(ctx context.Context, identity entity.Identity) error

	NewHTTP(r *mux.Router)
	HTTPLogin(w http.ResponseWriter, r *http.Request)
	HTTPRegister(w http.ResponseWriter, r *http.Request)
}

type service struct {
	db              *gorm.DB
	secret          string
	tokenExpiration int
	store           *sessions.CookieStore
	enforce         *casbin.Enforcer
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB, secret string, tokenExpiration int, store *sessions.CookieStore, enforce *casbin.Enforcer) Service {
	return service{db, secret, tokenExpiration, store, enforce}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, identity entity.Identity) (tokenAccess, tokenRefresh entity.Token, err error) {
	l := logger.Logger(ctx)

	identityDB := entity.Identity{}
	err = s.db.Model(&entity.Identity{}).Where("provider = ? AND uid = ?", identity.Provider, identity.UID).First(&identityDB).Error
	if err != nil {
		l.Warn().Caller().Msg(err.Error())
		return tokenAccess, tokenRefresh, e.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(identityDB.Password), []byte(identity.Password)); err != nil {
		l.Info().Caller().Msg("authentication failed")
		return tokenAccess, tokenRefresh, e.ErrUnauthorized
	}

	// Get user info.
	user := entity.User{}
	err = s.db.Model(&entity.User{}).Where("id = ?", identityDB.UserID).First(&user).Error
	if err != nil {
		l.Error().Caller().Msg(err.Error())
		return tokenAccess, tokenRefresh, err
	}

	tokenAccess, err = jwt.GenerateAccessToken(s.secret, identityDB, user)
	if err != nil {
		l.Error().Caller().Msg(err.Error())
		return tokenAccess, tokenRefresh, errors.New("error to generate access token: " + err.Error())
	}

	tokenRefresh, err = jwt.GenerateRefreshToken(s.secret, identityDB, user)
	if err != nil {
		l.Error().Caller().Msg(err.Error())
		return tokenAccess, tokenRefresh, errors.New("error to generate refresh token: " + err.Error())
	}

	tokenRefresh.UserID = identityDB.UserID
	err = s.db.Model(&entity.Token{}).Create(&tokenRefresh).Error
	if err != nil {
		l.Error().Caller().Msg(err.Error())
	}
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
