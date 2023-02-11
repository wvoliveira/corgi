package group

import (
	"database/sql"
	"errors"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Add(*gin.Context, model.Group, string) (model.Group, error)
	List(*gin.Context, int, int, string, string) (int64, int, []model.Group, error)
	FindByID(*gin.Context, string, string) (model.Group, []model.User, error)

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
	log := logger.Logger(c)

	sttCount, _ := s.db.PrepareContext(c, "SELECT COUNT(0) FROM group_user WHERE user_id = $1")

	err = sttCount.QueryRowContext(c, userID).Scan(&total)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	// TODO: fix to use "sort" variable
	sttData, _ := s.db.PrepareContext(c, `
		SELECT g.* FROM groups g
		JOIN group_user gu ON gu.user_id = $1
		GROUP BY g.id
		ORDER BY g.id ASC OFFSET $2 LIMIT $3
	`)

	rows, err := sttData.QueryContext(c, userID, offset, limit)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	defer rows.Close()
	var group model.Group

	for rows.Next() {
		err = rows.Scan(&group.ID, &group.CreatedAt, &group.UpdatedAt, &group.Name, &group.DisplayName, &group.Description, &group.CreatedBy, &group.OwnerID)

		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return
		}

		groups = append(groups, group)
	}

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return total, pages, groups, e.ErrInternalServerError
	}

	pages = int(math.Ceil(float64(total) / float64(limit)))

	return
}

// FindByID get a group details filtering with group ID.
func (s service) FindByID(c *gin.Context, groupID, userID string) (group model.Group, users []model.User, err error) {
	log := logger.Logger(c)

	// Get group details.
	sttGroup, _ := s.db.PrepareContext(c, `
		SELECT g.* FROM groups g
		JOIN group_user gu ON gu.user_id = $1
		WHERE g.id = $2
		GROUP BY g.id;
	`)

	err = sttGroup.QueryRowContext(c, userID, groupID).Scan(
		&group.ID,
		&group.CreatedAt,
		&group.UpdatedAt,
		&group.Name,
		&group.DisplayName,
		&group.Description,
		&group.CreatedBy,
		&group.OwnerID,
	)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	// Get all users by group ID.
	sttUsers, _ := s.db.PrepareContext(c, `
		SELECT u.id, u.name FROM users u
		JOIN groups g ON g.id = $1
		JOIN group_user gu ON gu.user_id = u.id
		GROUP BY u.id;
	`)

	rows, err := sttUsers.QueryContext(c, groupID)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	defer rows.Close()
	var user model.User

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name)

		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return
		}

		users = append(users, user)
	}

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return group, users, e.ErrInternalServerError
	}

	return
}
