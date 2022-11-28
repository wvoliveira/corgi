package password

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Service encapsulates the authentication logic.
type Service interface {
	Login(*gin.Context, entity.Identity) (entity.User, error)
	Register(*gin.Context, entity.Identity) error

	NewHTTP(*gin.RouterGroup)
	HTTPLogin(c *gin.Context)
	HTTPRegister(c *gin.Context)
}

type service struct {
	// TODO: still use cache or remove?
	db    *gorm.DB
	cache *badger.DB
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB, cache *badger.DB) Service {
	return service{db, cache}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
func (s service) Login(c *gin.Context, identity entity.Identity) (user entity.User, err error) {
	var (
		log        = logger.Logger(c.Request.Context())
		identityDB = entity.Identity{}
	)

	err = s.db.Model(&entity.Identity{}).
		Where("provider = ? AND uid = ?", identity.Provider, identity.UID).
		First(&identityDB).Error

	if err != nil {
		log.Warn().Caller().Msg(err.Error())
		return user, e.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(identityDB.Password), []byte(identity.Password)); err != nil {
		log.Info().Caller().Msg("authentication failed")
		return user, e.ErrUnauthorized
	}

	err = s.db.Model(&entity.User{}).Where("id = ?", identityDB.UserID).First(&user).Error

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return user, err
	}

	return
}

// Register a new user to our database.
func (s service) Register(c *gin.Context, identity entity.Identity) (err error) {
	var (
		log        = logger.Logger(c.Request.Context())
		user       = entity.User{}
		identityDB = entity.Identity{}
	)

	err = s.db.Model(&entity.Identity{}).
		Where("provider = ? AND uid = ?", identity.Provider, identity.UID).
		First(&identityDB).Error

	if err == nil {
		log.Warn().Caller().Msg(fmt.Sprintf("provider '%s' and uid '%s' already exists", identity.Provider, identity.UID))
		return e.ErrAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(identity.Password), 8)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return e.ErrInternalServerError
	}

	identity.ID = uuid.New().String()
	identity.CreatedAt = time.Now()
	identity.Password = string(hashedPassword)

	active := true
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.Role = "user"
	user.Active = &active
	user.Identities = append(user.Identities, identity)

	err = s.db.Model(&entity.User{}).Create(&user).Error

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return errors.New("error to create a user: " + err.Error())
	}

	return
}
