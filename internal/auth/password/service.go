package password

import (
	"context"
	"errors"
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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
	signingKey      string
	tokenExpiration int
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB, signingKey string, tokenExpiration int, logger log.Logger) Service {
	return service{logger, db, signingKey, tokenExpiration}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, identity entity.Identity) (token entity.Token, err error) {
	logger := s.logger.With(ctx, identity.Provider, identity.UID)

	identityDB := entity.Identity{}
	err = s.db.Debug().Model(&entity.Identity{}).Where("provider = ? AND uid = ?", identity.Provider, identity.UID).First(&identityDB).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return token, err
	} else if err != nil {
		return token, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(identity.Password), []byte(identityDB.Password)); err != nil {
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

	accessToken, err := s.generateAccessToken(identityDB, user)
	if err != nil {
		return token, errors.New("error to generate access token: " + err.Error())
	}

	refreshToken, err := s.generateRefreshToken(identityDB, user)
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
	token.AtExpires = accessToken.AtExpires
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

func (s service) generateAccessToken(identity entity.Identity, user entity.User) (token entity.Token, err error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)

	// This system is not a security problem. So, the token expires in 2 hours.
	tokenExpires := time.Now().Add(time.Hour * 2).Unix()

	claims := accessToken.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["identity_id"] = identity.GetID()
	claims["identity_provider"] = identity.Provider // e-mail, google, facebook, etc.
	claims["identity_uid"] = identity.GetUID() // e-mail address, google id, facebook id, etc.
	claims["user_id"] = user.GetID()
	claims["user_role"] = user.GetRole()
	claims["exp"] = tokenExpires

	at, err := accessToken.SignedString([]byte(s.signingKey))
	if err != nil {
		err = errors.New("error to generate access token: " + err.Error())
		return
	}
	token.CreatedAt = time.Now()
	token.AccessToken = at
	token.AtExpires = tokenExpires
	token.UserID = user.GetID()
	return
}

func (s service) generateRefreshToken(identity entity.Identity, user entity.User) (token entity.Token, err error) {
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	id := uuid.New().String()

	// Refresh token expires in 7 days. But I think to increase this value.
	tokenExpires := time.Now().AddDate(0, 0, 7).Unix()

	claims := refreshToken.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["sub"] = 1
	claims["identity_id"] = identity.GetID()
	claims["identity_provider"] = identity.Provider // e-mail, google, facebook, etc.
	claims["identity_uid"] = identity.GetUID() // e-mail address, google id, facebook id, etc.
	claims["user_id"] = user.GetID()
	claims["user_role"] = user.GetRole()
	claims["exp"] = tokenExpires

	rt, err := refreshToken.SignedString([]byte(s.signingKey))
	if err != nil {
		err = errors.New("error to generate refresh token: " + err.Error())
		return
	}

	token.ID = id
	token.CreatedAt = time.Now()
	token.RefreshToken = rt
	token.RtExpires = tokenExpires
	token.UserID = user.GetID()
	return
}
