package password

import (
	"context"
	"errors"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/elga-io/corgi/internal/app/entity"
	e "github.com/elga-io/corgi/internal/pkg/errors"
	"github.com/elga-io/corgi/internal/pkg/jwt"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Service encapsulates the authentication logic.
type Service interface {
	Login(ctx context.Context, identity entity.Identity) (entity.Token, entity.Token, error)
	Register(ctx context.Context, identity entity.Identity) error

	NewHTTP(r *gin.Engine)
	HTTPLogin(c *gin.Context)
	HTTPRegister(c *gin.Context)
}

// Identity represents an authenticated user identity.
type Identity interface {
	// GetID returns the user ID.
	GetID() string
	// GetUID returns the e-mail, google id, facebook id, etc.
	GetUID() string
	// GetRole returns the role.
	GetRole() string
}

type service struct {
	db              *gorm.DB
	secret          string
	tokenExpiration int
	store           cookie.Store
	enforce         *casbin.Enforcer
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB, secret string, tokenExpiration int, store cookie.Store, enforce *casbin.Enforcer) Service {
	return service{db, secret, tokenExpiration, store, enforce}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, identity entity.Identity) (tokenAccess, tokenRefresh entity.Token, err error) {
	l := log.Ctx(ctx)

	identityDB := entity.Identity{}
	err = s.db.Model(&entity.Identity{}).Where("provider = ? AND uid = ?", identity.Provider, identity.UID).First(&identityDB).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		l.Warn().Caller().Msg(err.Error())
		return tokenAccess, tokenRefresh, e.ErrUnauthorized
	} else if err != nil {
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
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return tokenAccess, tokenRefresh, err
	} else if err != nil {
		return tokenAccess, tokenRefresh, err
	}

	tokenAccess, err = jwt.GenerateAccessToken(s.secret, identityDB, user)
	if err != nil {
		return tokenAccess, tokenRefresh, errors.New("error to generate access token: " + err.Error())
	}

	tokenRefresh, err = jwt.GenerateRefreshToken(s.secret, identityDB, user)
	if err != nil {
		return tokenAccess, tokenRefresh, errors.New("error to generate refresh token: " + err.Error())
	}

	tokenRefresh.UserID = identityDB.UserID
	err = s.db.Model(&entity.Token{}).Create(&tokenRefresh).Error
	if err != nil {
		return
	}
	return
}

// Register a new user to our database.
func (s service) Register(ctx context.Context, identity entity.Identity) (err error) {
	l := log.Ctx(ctx)

	user := entity.User{}
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
	return
}
