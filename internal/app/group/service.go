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
	Delete(*gin.Context, string, string) error
	InvitesAddByID(*gin.Context, invitesAddByIDRequest) (model.GroupInvite, error)
	InvitesListByID(*gin.Context, invitesListByIDRequest) (int64, int, []model.GroupInvite, error)
	InvitesList(*gin.Context, invitesListRequest) (int64, int, []model.GroupInvite, error)

	NewHTTP(*gin.RouterGroup)
	HTTPAdd(*gin.Context)
	HTTPList(*gin.Context)
	HTTPFindByID(*gin.Context)
	HTTPDelete(*gin.Context)
	HTTPInvitesAddByID(*gin.Context)
	HTTPInvitesListByID(*gin.Context)
	HTTPInvitesList(*gin.Context)
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
	group := model.Group{}

	for rows.Next() {
		err = rows.Scan(&group.ID, &group.CreatedAt, &group.UpdatedAtNull, &group.Name, &group.DisplayName, &group.Description, &group.CreatedBy, &group.OwnerID)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return
		}

		if group.UpdatedAtNull.Valid {
			group.UpdatedAt = &group.UpdatedAtNull.Time
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
func (s service) FindByID(c *gin.Context, whoID, groupID string) (group model.Group, users []model.User, err error) {
	log := logger.Logger(c)

	query := `
		SELECT g.* FROM groups g
		INNER JOIN group_user gu ON gu.group_id = g.id
		INNER JOIN users u ON u.id = gu.user_id
		WHERE u.id = $1
		AND g.id = $2
	`

	log.Debug().Caller().Msg(query)
	sttGroup, _ := s.db.PrepareContext(c, query)

	err = sttGroup.QueryRowContext(c, whoID, groupID).Scan(
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

	query = `
		SELECT u.id, u.name FROM users u
		INNER JOIN group_user gu ON gu.user_id = u.id
		INNER JOIN groups g ON g.id = gu.group_id
		WHERE u.id = $1
		AND g.id = $2
	`

	// Get all users by group ID.
	log.Debug().Caller().Msg(query)
	sttUsers, _ := s.db.PrepareContext(c, query)

	rows, err := sttUsers.QueryContext(c, whoID, groupID)
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
func (s service) Delete(c *gin.Context, whoID, groupID string) (err error) {
	log := logger.Logger(c)

	group := model.Group{}

	query := `SELECT id FROM groups WHERE id = $1 AND owner_id = $2`
	log.Debug().Caller().Msg(query)

	err = s.db.QueryRowContext(c, query, groupID, whoID).Scan(&group.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error().Caller().Msg(err.Error())
			return e.ErrInternalServerError
		}
	}

	if group.ID == "" {
		message := fmt.Sprintf("Group with id = '%s' and owner_id = '%s' was not found", groupID, whoID)
		log.Info().Caller().Msg(message)
		return e.ErrGroupNotFound
	}

	key := fmt.Sprintf("group_%s", groupID)

	// Keep going on error from cache.
	_, err = s.cache.Del(c, key).Result()
	if err != nil {
		log.Error().Caller().Msg(err.Error())
	}

	query = "DELETE FROM groups WHERE id = $1 AND owner_id = $2"

	_, err = s.db.ExecContext(c, query, groupID, whoID)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return e.ErrInternalServerError
	}

	return
}

func (s service) InvitesAddByID(c *gin.Context, payload invitesAddByIDRequest) (groupInvite model.GroupInvite, err error) {
	log := logger.Logger(c)

	// Check if invite already exists.
	query := `
		SELECT gi.id FROM groups_invites gi
		INNER JOIN identities i ON i.user_id = gi.user_id
		WHERE gi.group_id = $1
		AND i.provider = 'email'
		AND i.uid = $2
	`
	log.Debug().Caller().Msg(query)

	stt, err := s.db.PrepareContext(c, query)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	err = stt.QueryRowContext(c, payload.GroupID, payload.UserEmail).Scan(&groupInvite.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Error().Caller().Msg(err.Error())
			return
		}
	}

	// Get user ID from identities.
	query = `SELECT user_id FROM identities WHERE provider = 'email' AND uid = $1`
	log.Debug().Caller().Msg(query)

	stt, err = s.db.PrepareContext(c, query)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	user := model.User{}
	err = stt.QueryRowContext(c, payload.UserEmail).Scan(&user.ID)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	if user.ID == "" {
		err = e.ErrUserNotFound
		return
	}

	groupInvite.ID = ulid.Make().String()
	query = `INSERT INTO groups_invites(id, group_id, user_id, invited_by) VALUES($1, $2, $3, $4)`
	log.Debug().Caller().Msg(query)

	stt, err = s.db.PrepareContext(c, query)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	_, err = stt.ExecContext(c, groupInvite.ID, payload.GroupID, user.ID, payload.InvitedBy)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	return
}

func (s service) InvitesListByID(c *gin.Context, payload invitesListByIDRequest) (total int64, pages int, invites []model.GroupInvite, err error) {
	log := logger.Logger(c)

	sttCount, _ := s.db.PrepareContext(c, "SELECT COUNT(0) FROM group_user WHERE group_id = $1")
	err = sttCount.QueryRowContext(c, payload.GroupID).Scan(&total)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	// TODO: fix to use "sort" variable
	query := `SELECT gi.* FROM groups_invites gi
		INNER JOIN groups g on g.id = gi.group_id
		WHERE g.id = $1
		ORDER BY g.id ASC OFFSET $2 LIMIT $3
	`

	log.Debug().Caller().Msg(query)
	sttData, _ := s.db.PrepareContext(c, query)

	rows, err := sttData.QueryContext(c, payload.GroupID, payload.Offset, payload.Limit)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	defer rows.Close()
	invites = []model.GroupInvite{}
	gi := model.GroupInvite{}

	for rows.Next() {
		err = rows.Scan(&gi.ID, &gi.CreatedAt, &gi.UpdatedAtNull, &gi.GroupID, &gi.UserID, &gi.InvitedBy, &gi.Accepted)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return
		}

		if gi.UpdatedAtNull.Valid {
			gi.UpdatedAt = &gi.UpdatedAtNull.Time
		}

		invites = append(invites, gi)
	}

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return total, pages, invites, e.ErrInternalServerError
	}

	pages = int(math.Ceil(float64(total) / float64(payload.Limit)))
	return
}

func (s service) InvitesList(c *gin.Context, payload invitesListRequest) (total int64, pages int, invites []model.GroupInvite, err error) {
	log := logger.Logger(c)

	query := `
		SELECT COUNT(0) FROM groups_invites
		INNER JOIN group_user gu on gu.group_id = groups_invites.group_id
		WHERE gu.user_id = $1
	`

	sttCount, _ := s.db.PrepareContext(c, query)
	err = sttCount.QueryRowContext(c, payload.WhoID).Scan(&total)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	// TODO: fix to use "sort" variable
	query = `
		SELECT gi.* FROM groups_invites gi
		INNER JOIN group_user gu on gu.group_id = gi.group_id
		WHERE gu.user_id = $1
		ORDER BY gi.id ASC OFFSET $2 LIMIT $3
	`

	log.Debug().Caller().Msg(query)
	sttData, _ := s.db.PrepareContext(c, query)

	rows, err := sttData.QueryContext(c, payload.WhoID, payload.Offset, payload.Limit)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	defer rows.Close()
	invites = []model.GroupInvite{}
	gi := model.GroupInvite{}

	for rows.Next() {
		err = rows.Scan(&gi.ID, &gi.CreatedAt, &gi.UpdatedAtNull, &gi.GroupID, &gi.UserID, &gi.InvitedBy, &gi.Accepted)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return
		}

		if gi.UpdatedAtNull.Valid {
			gi.UpdatedAt = &gi.UpdatedAtNull.Time
		}

		invites = append(invites, gi)
	}

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return total, pages, invites, e.ErrInternalServerError
	}

	pages = int(math.Ceil(float64(total) / float64(payload.Limit)))
	return
}
