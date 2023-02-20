package group

import (
	"database/sql"
	"errors"
	"fmt"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Add(*gin.Context, string, model.Group) (model.Group, error)
	List(*gin.Context, string, int, int, string) (int64, int, []model.Group, error)
	FindByID(*gin.Context, string, string) (model.Group, []model.User, error)

	InviteAdd(*gin.Context, model.GroupInvite) (model.GroupInvite, error)

	NewHTTP(*gin.RouterGroup)
	HTTPAdd(*gin.Context)
	HTTPList(*gin.Context)
}

type service struct {
	db    *sql.DB
	cache *redis.Client
}

// NewService creates a new group service.
func NewService(db *sql.DB, cache *redis.Client) Service {
	return service{db, cache}
}

func (s service) Add(c *gin.Context, whoID string, payload model.Group) (group model.Group, err error) {
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
	group.CreatedBy = whoID
	group.OwnerID = whoID

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

	_, err = stmt.ExecContext(c, group.ID, whoID)
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

func (s service) List(c *gin.Context, whoID string, offset, limit int, sort string) (total int64, pages int, groups []model.Group, err error) {
	log := logger.Logger(c)

	sttCount, _ := s.db.PrepareContext(c, "SELECT COUNT(0) FROM group_user WHERE user_id = $1")

	err = sttCount.QueryRowContext(c, whoID).Scan(&total)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	// TODO: fix to use "sort" variable
	query := `SELECT g.* FROM groups g
		INNER JOIN group_user gu ON gu.group_id = g.id
		INNER JOIN users u ON u.id = gu.user_id
		WHERE u.id = $1
		ORDER BY g.id ASC OFFSET $2 LIMIT $3`

	log.Debug().Caller().Msg(query)
	sttData, _ := s.db.PrepareContext(c, query)

	rows, err := sttData.QueryContext(c, whoID, offset, limit)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	defer rows.Close()

	groups = []model.Group{}
	var group model.Group

	for rows.Next() {
		err = rows.Scan(&group.ID, &group.CreatedAt, &group.UpdatedAtNull, &group.Name, &group.DisplayName, &group.Description, &group.CreatedBy, &group.OwnerID)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return
		}

		fmt.Println(group.UpdatedAtNull)

		if group.UpdatedAtNull.Valid {
			group.UpdatedAt = group.UpdatedAtNull.Time
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
		INNER JOIN group_user gu ON gu.group_id = g.id
		INNER JOIN users u ON u.id = gu.user_id
		WHERE u.id = $1
		AND g.id = $2
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
		if errors.Is(err, sql.ErrNoRows) {
			err = e.ErrGroupNotFound
			return
		}

		log.Error().Caller().Msg(err.Error())
		return
	}

	// Get all users by group ID.
	sttUsers, _ := s.db.PrepareContext(c, `
		SELECT u.id, u.name FROM users u
		INNER JOIN group_user gu ON gu.user_id = u.id
		INNER JOIN groups g ON g.id = gu.group_id
		WHERE u.id = $1
		AND g.id = $2
	`)

	rows, err := sttUsers.QueryContext(c, userID, groupID)

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

// Delete delete a group by ID.
func (s service) Delete(c *gin.Context, userID, groupID string) (err error) {
	log := logger.Logger(c)

	group := model.Group{}

	query := "SELECT id FROM groups WHERE id = $1 AND owner_id = $2"
	err = s.db.QueryRowContext(c, query, groupID, userID).Scan(&group.ID)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Error().Caller().Msg(err.Error())
			return e.ErrInternalServerError
		}
	}

	if group.ID == "" {
		message := fmt.Sprintf("Group with id = '%s' and owner_id = '%s' was not found", groupID, userID)
		log.Info().Caller().Msg(message)
		return e.ErrGroupNotFound
	}

	key := fmt.Sprintf("group_%s", groupID)
	_, err = s.cache.Del(c, key).Result()

	// Keep going on error from cache.
	if err != nil {
		log.Error().Caller().Msg(err.Error())
	}

	query = "DELETE FROM groups WHERE id = $1 AND owner_id = $2"
	_, err = s.db.ExecContext(c, query, groupID, userID)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return e.ErrInternalServerError
	}

	return
}

func (s service) InviteAdd(c *gin.Context, payload model.GroupInvite) (groupInvite model.GroupInvite, err error) {
	log := logger.Logger(c)

	stt, err := s.db.PrepareContext(c, "SELECT id FROM groups_invites WHERE group_id = $1 AND user_id = $2 AND invited_by = $3")
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	err = stt.QueryRowContext(c, payload.GroupID, payload.UserID, payload.InvitedBy).Scan(&groupInvite.ID)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Error().Caller().Msg(err.Error())
			return
		}
	}

	user := model.User{}

	stt, err = s.db.PrepareContext(c, "SELECT id FROM users WHERE NOT id = $1 AND id = $2")
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	err = stt.QueryRowContext(c, 0, payload.UserID).Scan(&user.ID)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	if user.ID == "" {
		err = e.ErrUserNotFound
		return
	}

	payload.ID = ulid.Make().String()

	stt, err = s.db.PrepareContext(c, "INSERT INTO groups_invites(id, group_id, user_id, invited_by) VALUES($1, $2, $3, $4)")
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	_, err = stt.ExecContext(c, payload.ID, payload.GroupID, payload.UserID, payload.InvitedBy)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	return
}
