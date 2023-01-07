package group

import (
	"errors"
	"fmt"
	"math"

	"github.com/gin-gonic/gin"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	"gorm.io/gorm"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Add(*gin.Context, model.Group, string) (model.Group, error)
	List(*gin.Context, int, int, string, string) (int64, int, []model.Group, error)
	FindByID(*gin.Context, string, string) (model.Group, error)

	NewHTTP(*gin.RouterGroup)
	HTTPAdd(*gin.Context)
	HTTPList(*gin.Context)
}

type service struct {
	db *gorm.DB
}

// NewService creates a new group service.
func NewService(database *gorm.DB) Service {
	return service{database}
}

func (s service) Add(c *gin.Context, requestGroup model.Group, userID string) (group model.Group, err error) {
	l := logger.Logger(c)

	err = s.db.Model(&model.Group{}).
		Where("name = ?", requestGroup.Name).
		Take(&group).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		group = requestGroup
		group.CreatedBy = userID
		group.Users = append(group.Users, model.User{ID: userID})

		err = s.db.Create(&group).Save(&group).Error
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

func (s service) List(c *gin.Context, offset, limit int, sort, userID string) (total int64, pages int, groups []model.Group, err error) {
	l := logger.Logger(c)
	pages = 1

	// TODO: make a single transaction to get total and list of items.
	// For some reason it's a bit complex to make join with gorm
	// or my brain is not ready for this.
	queryTotal := `SELECT COUNT()
		FROM groups 
		JOIN user_groups 
			ON user_groups.group_id = groups.id 
			AND user_groups.user_id = ?
	`
	err = s.db.Raw(queryTotal, userID).Scan(&total).Error
	if err == gorm.ErrRecordNotFound {
		return
	}

	if err != nil {
		l.Error().Caller().Msg(err.Error())
		return
	}

	queryItems := fmt.Sprintf(`SELECT groups.* 
		FROM groups 
		JOIN user_groups 
			ON user_groups.group_id = groups.id 
			AND user_groups.user_id = ?
		ORDER BY %s
		LIMIT ? OFFSET ?
	`, sort)
	err = s.db.Raw(queryItems, userID, limit, offset).Scan(&groups).Error
	if err == gorm.ErrRecordNotFound {
		return
	}

	if err != nil {
		l.Error().Caller().Msg(err.Error())
		return
	}

	pages = int(math.Ceil(float64(total) / float64(limit)))
	return
}

// FindByID get a group details filtering with group ID.
func (s service) FindByID(c *gin.Context, groupID, userID string) (group model.Group, err error) {
	l := logger.Logger(c)

	query := `SELECT groups.*
		FROM groups 
		JOIN user_groups 
			ON user_groups.group_id = groups.id 
			AND user_groups.user_id = ?
		WHERE groups.id = ?
	`
	err = s.db.Raw(query, userID, groupID).Scan(&group).Error

	if err == gorm.ErrRecordNotFound || group.ID == "" {
		return group, e.ErrGroupNotFound
	}

	if err != nil {
		l.Error().Caller().Msg(err.Error())
	}

	users := []model.User{}
	queryUsers := `SELECT users.*
		FROM users 
		JOIN user_groups 
			ON user_groups.user_id = users.id 
		JOIN groups
			ON groups.id = user_groups.group_id 
		AND groups.id = ?
	`
	err = s.db.Raw(queryUsers, groupID).Scan(&users).Error

	if err == gorm.ErrRecordNotFound {
		l.Warn().Caller().Msg(err.Error())
		return group, e.ErrGroupNotFound
	}

	if err != nil {
		l.Error().Caller().Msg(err.Error())
	}

	group.Users = users
	return
}
