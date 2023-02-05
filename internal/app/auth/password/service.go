package password

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	"golang.org/x/crypto/bcrypt"
)

// Service encapsulates the authentication logic.
type Service interface {
	Login(*gin.Context, model.Identity) (model.User, error)
	Register(*gin.Context, model.Identity) error

	NewHTTP(*gin.RouterGroup)
	HTTPLogin(c *gin.Context)
	HTTPRegister(c *gin.Context)
}

type service struct {
	// TODO: still use cache or remove?
	db *sql.DB
	kv *badger.DB
}

// NewService creates a new authentication service.
func NewService(db *sql.DB, kv *badger.DB) Service {
	return service{db, kv}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
func (s service) Login(c *gin.Context, identity model.Identity) (user model.User, err error) {
	var (
		log        = logger.Logger(c.Request.Context())
		identityDB = model.Identity{}
	)

	err = s.db.Model(&model.Identity{}).
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

	err = s.db.Model(&model.User{}).Where("id = ?", identityDB.UserID).First(&user).Error

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return user, e.ErrAuthPasswordInternalError
	}

	return
}

// Register a new user to our database.
func (s service) Register(c *gin.Context, identity model.Identity) (err error) {
	var (
		log        = logger.Logger(c.Request.Context())
		user       = model.User{}
		identityDB = model.Identity{}
	)

	err = s.db.Model(&model.Identity{}).
		Where("provider = ? AND uid = ?", identity.Provider, identity.UID).
		First(&identityDB).Error

	if err == nil {
		log.Warn().Caller().Msg(fmt.Sprintf("provider '%s' and uid '%s' already exists", identity.Provider, identity.UID))
		return e.ErrAuthPasswordUserAlreadyExists
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

	err = s.db.Model(&model.User{}).Create(&user).Error

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return e.ErrAuthPasswordInternalError
	}

	return
}
