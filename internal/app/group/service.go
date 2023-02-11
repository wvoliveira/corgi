package group

import (
	"database/sql"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
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
	db *sql.DB
}

// NewService creates a new group service.
func NewService(database *sql.DB) Service {
	return service{database}
}

func (s service) Add(c *gin.Context, payload model.Group, userID string) (group model.Group, err error) {
	log := logger.Logger(c)

	err = s.db.QueryRowContext(c, "SELECT id FROM groups WHERE name = $1", payload.Name).Scan(&group)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Error().Caller().Msg(err.Error())
			return
		}
	}

	group = payload
	group.ID = ulid.Make().String()
	group.CreatedBy = userID
	group.OwnerID = userID

	tx, err := s.db.BeginTx(c, &sql.TxOptions{})
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	// TODO:
	// 	- create a relation many to many and fix this query.
	// 	- check error when rollback
	stmt, err := tx.PrepareContext(c,
		"INSERT INTO groups(id, name, display_name, description, created_by, owner_id) VALUES($1, $2, $3, $4, $5, $6)",
	)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	_, err = stmt.ExecContext(c, group.ID, group.Name, group.DisplayName, group.Description, group.CreatedBy, group.OwnerID)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		_ = tx.Rollback()
		return
	}

	stmt, err = tx.PrepareContext(c, "INSERT INTO group_user(group_id, user_id) VALUES($1, $2)")
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	_, err = stmt.ExecContext(c, group.ID, userID)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		_ = tx.Rollback()
		return
	}

	err = tx.Commit()

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		_ = tx.Rollback()
		return
	}

	return
}

func (s service) List(c *gin.Context, offset, limit int, sort, userID string) (total int64, pages int, groups []model.Group, err error) {
	// l := logger.Logger(c)
	// pages = 1

	// // TODO: make a single transaction to get total and list of items.
	// // For some reason it's a bit complex to make join with gorm
	// // or my brain is not ready for this.
	// queryTotal := `SELECT COUNT()
	// 	FROM groups
	// 	JOIN user_groups
	// 		ON user_groups.group_id = groups.id
	// 		AND user_groups.user_id = ?
	// `
	// err = s.db.Raw(queryTotal, userID).Scan(&total).Error
	// if err == gorm.ErrRecordNotFound {
	// 	return
	// }

	// if err != nil {
	// 	l.Error().Caller().Msg(err.Error())
	// 	return
	// }

	// queryItems := fmt.Sprintf(`SELECT groups.*
	// 	FROM groups
	// 	JOIN user_groups
	// 		ON user_groups.group_id = groups.id
	// 		AND user_groups.user_id = ?
	// 	ORDER BY %s
	// 	LIMIT ? OFFSET ?
	// `, sort)
	// err = s.db.Raw(queryItems, userID, limit, offset).Scan(&groups).Error
	// if err == gorm.ErrRecordNotFound {
	// 	return
	// }

	// if err != nil {
	// 	l.Error().Caller().Msg(err.Error())
	// 	return
	// }

	// pages = int(math.Ceil(float64(total) / float64(limit)))
	return
}

// FindByID get a group details filtering with group ID.
func (s service) FindByID(c *gin.Context, groupID, userID string) (group model.Group, err error) {
	// l := logger.Logger(c)

	// query := `SELECT groups.*
	// 	FROM groups
	// 	JOIN user_groups
	// 		ON user_groups.group_id = groups.id
	// 		AND user_groups.user_id = ?
	// 	WHERE groups.id = ?
	// `
	// err = s.db.Raw(query, userID, groupID).Scan(&group).Error

	// if err == gorm.ErrRecordNotFound || group.ID == "" {
	// 	return group, e.ErrGroupNotFound
	// }

	// if err != nil {
	// 	l.Error().Caller().Msg(err.Error())
	// }

	// users := []model.User{}
	// queryUsers := `SELECT users.*
	// 	FROM users
	// 	JOIN user_groups
	// 		ON user_groups.user_id = users.id
	// 	JOIN groups
	// 		ON groups.id = user_groups.group_id
	// 	AND groups.id = ?
	// `
	// err = s.db.Raw(queryUsers, groupID).Scan(&users).Error

	// if err == gorm.ErrRecordNotFound {
	// 	l.Warn().Caller().Msg(err.Error())
	// 	return group, e.ErrGroupNotFound
	// }

	// if err != nil {
	// 	l.Error().Caller().Msg(err.Error())
	// }

	// group.Users = users
	return
}
