package group

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"gorm.io/gorm"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Add(ctx context.Context, requestGroup entity.Group, userID string) (group entity.Group, err error)
	List(ctx context.Context, offset, limit int, sort, userID string) (total int64, pages int, groups []entity.Group, err error)

	NewHTTP(r *mux.Router)
	HTTPAdd(w http.ResponseWriter, r *http.Request)
	HTTPList(w http.ResponseWriter, r *http.Request)
}

type service struct {
	db     *gorm.DB
	secret string
}

// NewService creates a new group service.
func NewService(database *gorm.DB, secret string) Service {
	return service{database, secret}
}

func (s service) Add(ctx context.Context, requestGroup entity.Group, userID string) (group entity.Group, err error) {
	l := logger.Logger(ctx)

	err = s.db.Model(&entity.Group{}).
		Where("name = ?", requestGroup.Name).
		Take(&group).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		group = requestGroup
		err = s.db.Model(&entity.Group{}).Create(&group).Error
		if err == nil {
			l.Info().Caller().Msg("group created with successfully")
		}
		return

	} else if err == nil {
		l.Warn().Caller().Msg(fmt.Sprintf("group with '%s' name already exists. Choose another one", group.Name))
		return group, e.ErrGroupAlreadyExists
	}

	l.Error().Caller().Msg(err.Error())
	return group, e.ErrInternalServerError
}

func (s service) List(ctx context.Context, offset, limit int, sort, userID string) (total int64, pages int, groups []entity.Group, err error) {
	l := logger.Logger(ctx)

	err = s.db.Model(&entity.Group{}).Where("user_id = ?", userID).
		Count(&total).
		Offset(offset).
		Limit(limit).
		Order(sort).
		Find(&groups).Error

	if err == gorm.ErrRecordNotFound || err == nil {
		l.Info().Caller().Msg("groups empty")
		return
	}

	pages = int(math.Ceil(float64(total) / float64(limit)))

	l.Error().Caller().Msg(err.Error())
	return
}
