package password

import (
	"context"
	"errors"
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/jwt"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

// Service encapsulates the authentication logic.
type Service interface {
	Login(ctx context.Context, identity entity.Identity) (entity.Token, error)
	Register(ctx context.Context, identity entity.Identity) error

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
	logger          log.Logger
	db              *gorm.DB
	secret      string
	tokenExpiration int
}

// NewService creates a new authentication service.
func NewService(logger log.Logger, db *gorm.DB, secret string, tokenExpiration int) Service {
	return service{logger, db, secret, tokenExpiration}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, identity entity.Identity) (token entity.Token, err error) {
	logger := s.logger.With(ctx, identity.Provider, identity.UID)

	identityDB := entity.Identity{}
	err = s.db.Debug().Model(&entity.Identity{}).Where("provider = ? AND uid = ?", identity.Provider, identity.UID).First(&identityDB).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("this provider not found in database")
		return token, err
	} else if err != nil {
		logger.Error("error when get identity from database", err.Error())
		return token, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(identityDB.Password), []byte(identity.Password)); err != nil {
		logger.Infof("authentication failed")
		return token, e.ErrUnauthorized
	}

	// Get user info.
	user := entity.User{}
	err = s.db.Debug().Model(&entity.User{}).Where("id = ?", identityDB.UserID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return token, err
	} else if err != nil {
		return token, err
	}

	accessToken, err := jwt.GenerateAccessToken(s.secret, identityDB, user)
	if err != nil {
		return token, errors.New("error to generate access token: " + err.Error())
	}

	refreshToken, err := jwt.GenerateRefreshToken(s.secret, identityDB, user)
	if err != nil {
		return token, errors.New("error to generate refresh token: " + err.Error())
	}

	refreshToken.AccessToken = accessToken.AccessToken
	refreshToken.UserID = identityDB.UserID
	err = s.db.Debug().Model(&entity.Token{}).Create(&refreshToken).Error
	if err != nil {
		return
	}

	token.AccessToken = accessToken.AccessToken
	token.RefreshToken = refreshToken.RefreshToken
	token.AccessExpires = accessToken.AccessExpires
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

	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.Role = "user"
	user.Active = "true"
	user.Identities = append(user.Identities, identity)

	err = s.db.Debug().Model(&entity.User{}).Create(&user).Error
	return
}
