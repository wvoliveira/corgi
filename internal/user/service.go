package user

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Find(ctx context.Context, userID string) (entity.User, error)
	Update(ctx context.Context, user entity.User) (entity.User, error)

	HTTPFind(c *gin.Context)
	HTTPUpdate(c *gin.Context)

	Routers(r *gin.Engine)
}

type service struct {
	logger  log.Logger
	db      *gorm.DB
	secret  string
	store   cookie.Store
	enforce *casbin.Enforcer
}

// NewService creates a new authentication service.
func NewService(logger log.Logger, db *gorm.DB, secret string, store cookie.Store, enforce *casbin.Enforcer) Service {
	return service{logger, db, secret, store, enforce}
}

// Find get a shortener link from ID.
func (s service) Find(ctx context.Context, userID string) (user entity.User, err error) {
	logger := s.logger.With(ctx, "user_id", userID)

	user.ID = userID
	fmt.Println("USER_ID:", userID)

	var identities []entity.Identity
	err = s.db.Debug().Model(&user).Preload("Identities").Find(&user).Error
	if err == gorm.ErrRecordNotFound {
		logger.Infof("the user with user_id '%s' not found", userID)
		return user, e.ErrUserNotFound
	} else if err == nil {
		fmt.Println("1 user")
		fmt.Println(user)

		fmt.Println("1 identities")
		fmt.Println(identities)
		return
	}

	fmt.Println("2 identities")
	fmt.Println(identities)
	logger.Errorf("oh crap, an errors occurred: %s", err.Error())
	return
}

// Update change specific link by ID.
func (s service) Update(ctx context.Context, req entity.User) (user entity.User, err error) {
	logger := s.logger.With(ctx, "user_id", req.ID)
	logger.Infof("updating user with id '%s'", req.ID)

	req.UpdatedAt = time.Now()

	err = s.db.Model(&entity.User{}).Where("id = ?", req.ID).Updates(&req).Error
	if err == gorm.ErrRecordNotFound {
		logger.Infof("the user with id '%s' not found", req.ID)
		return user, e.ErrUserNotFound
	} else if err == nil {
		user = req
		return
	}

	logger.Errorf("oh crap, an errors occurred: %s", err.Error())
	return
}
