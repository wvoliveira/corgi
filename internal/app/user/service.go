package user

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"gorm.io/gorm"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Find(*gin.Context, string) (entity.User, error)
	Update(*gin.Context, entity.User) (entity.User, error)

	NewHTTP(*gin.RouterGroup)
	HTTPFind(*gin.Context)
	HTTPUpdate(*gin.Context)
}

type service struct {
	db    *gorm.DB
	cache *badger.DB
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB, cache *badger.DB) Service {
	return service{db, cache}
}

// Find get a shortener link from ID.
func (s service) Find(c *gin.Context, userID string) (user entity.User, err error) {
	var (
		log = logger.Logger(c.Request.Context())
	)

	if userID == "anonymous" {
		user.Name = "Anonymous"
		return
	}

	user.ID = userID

	err = s.db.Debug().
		Model(&user).
		Preload("Identities").
		Find(&user).Error

	if err == gorm.ErrRecordNotFound {
		log.Info().Caller().Msg(fmt.Sprintf("the user with user_id \"%s\" was not found", userID))
		return user, e.ErrUserNotFound
	}

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	return
}

// Update change specific link by ID.
func (s service) Update(c *gin.Context, reqUser entity.User) (user entity.User, err error) {
	var (
		log = logger.Logger(c.Request.Context())
	)

	if reqUser.ID == "anonymous" {
		return user, e.ErrUnauthorized
	}

	reqUser.UpdatedAt = time.Now()
	err = s.db.Model(&entity.User{}).
		Where("id = ?", reqUser.ID).
		Updates(&reqUser).Error

	if err == gorm.ErrRecordNotFound {
		log.Info().Caller().Msg(fmt.Sprintf("the user with id '%s' not found", reqUser.ID))
		return user, e.ErrUserNotFound
	}

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	user = reqUser
	return
}
