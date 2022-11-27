package password

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/jwt"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Service encapsulates the authentication logic.
type Service interface {
	Login(*gin.Context, entity.Identity) (entity.Token, entity.Token, error)
	Register(*gin.Context, entity.Identity) error

	NewHTTP(*gin.RouterGroup)
	HTTPLogin() gin.HandlerFunc
	HTTPRegister() gin.HandlerFunc
}

type service struct {
	db *gorm.DB
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB) Service {
	return service{db}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(c *gin.Context, identity entity.Identity) (tokenAccess, tokenRefresh entity.Token, err error) {
	l := logger.Logger(c.Request.Context())

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
func (s service) Register(c *gin.Context, identity entity.Identity) (err error) {
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
