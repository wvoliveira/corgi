package password

import (
	"context"
	"errors"
	"github.com/casbin/casbin/v2"
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/jwt"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

// Service encapsulates the authentication logic.
type Service interface {
	Login(ctx context.Context, identity entity.Identity) (entity.Token, entity.Token, error)
	Register(ctx context.Context, identity entity.Identity) error
	HTTPLogin(c *gin.Context)
	HTTPRegister(c *gin.Context)
	Routers(r *gin.Engine)
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
	logger          log.Logger
	db              *gorm.DB
	secret          string
	tokenExpiration int
	store           cookie.Store
	enforce         *casbin.Enforcer
}

// NewService creates a new authentication service.
func NewService(logger log.Logger, db *gorm.DB, secret string, tokenExpiration int, store cookie.Store, enforce *casbin.Enforcer) Service {
	return service{logger, db, secret, tokenExpiration, store, enforce}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, identity entity.Identity) (tokenAccess, tokenRefresh entity.Token, err error) {
	logger := s.logger.With(ctx, identity.Provider, identity.UID)

	identityDB := entity.Identity{}
	err = s.db.Model(&entity.Identity{}).Where("provider = ? AND uid = ?", identity.Provider, identity.UID).First(&identityDB).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Warnf("this provider + uid not found in database: %s", err.Error())
		return tokenAccess, tokenRefresh, e.ErrUnauthorized
	} else if err != nil {
		logger.Warnf("error when get identity from database: %s", err.Error())
		return tokenAccess, tokenRefresh, e.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(identityDB.Password), []byte(identity.Password)); err != nil {
		logger.Infof("authentication failed")
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
	logger := s.logger.With(ctx, identity.Provider, identity.UID)

	user := entity.User{}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(identity.Password), 8)
	if err != nil {
		logger.Error("error to create a hashed password:", err.Error())
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
