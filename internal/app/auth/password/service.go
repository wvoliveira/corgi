package password

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	"github.com/wvoliveira/corgi/internal/pkg/token"
)

// Service encapsulates the authentication logic.
type Service interface {
	Login(*gin.Context, model.Identity) (string, string, model.User, error)
	Register(*gin.Context, model.Identity, model.User) error

	NewHTTP(*gin.RouterGroup)
	HTTPLogin(c *gin.Context)
	HTTPRegister(c *gin.Context)
}

type service struct {
	// TODO: still use cache or remove?
	db    *sql.DB
	cache *redis.Client
}

// NewService creates a new authentication service.
func NewService(db *sql.DB, cache *redis.Client) Service {
	return service{db, cache}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
func (s service) Login(c *gin.Context, identity model.Identity) (accessToken, refreshToken string, user model.User, err error) {
	log := logger.Logger(c.Request.Context())
	identityFromDB := model.Identity{}

	query := "SELECT user_id, password FROM identities WHERE provider IN ('username', 'email') AND uid = $1"
	err = s.db.QueryRowContext(c, query, identity.UID).Scan(&identityFromDB.UserID, &identityFromDB.Password)

	if err != nil {
		log.Warn().Caller().Msg(err.Error())
		return accessToken, refreshToken, user, e.ErrUnauthorized
	}

	err = identity.CheckPassword(identityFromDB.Password)
	if err != nil {
		log.Info().Caller().Msg(fmt.Sprintf("username/email and password dont match: %s", err.Error()))
		return accessToken, refreshToken, user, e.ErrUnauthorized
	}

	query = "SELECT * FROM users WHERE id = $1"
	err = s.db.QueryRowContext(c, query, identityFromDB.UserID).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Role, &user.Active)

	if err != nil {
		log.Warn().Caller().Msg(err.Error())
		return accessToken, refreshToken, user, e.ErrAuthPasswordInternalError
	}

	accessToken, refreshToken, err = token.GenerateJWTAccess(user)
	return
}

// Register a new user to our database.
func (s service) Register(c *gin.Context, identity model.Identity, user model.User) (err error) {
	log := logger.Logger(c)

	idenDB := model.Identity{}

	query := "SELECT id FROM identities WHERE provider = $1 AND uid = $2"
	err = s.db.QueryRowContext(c, query, identity.Provider, identity.UID).Scan(&idenDB.ID)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Warn().Caller().Msg(err.Error())
			return
		}
	}

	if idenDB.ID != "" {
		message := fmt.Sprintf("provider '%s' and uid '%s' already exists", identity.Provider, identity.UID)
		log.Warn().Caller().Msg(message)
		return e.ErrAuthPasswordUserAlreadyExists
	}

	identity.ID = ulid.Make().String()
	identity.HashPassword(identity.Password)

	user.ID = ulid.Make().String()
	user.Role = "user"

	tx, err := s.db.BeginTx(c, &sql.TxOptions{})
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return e.ErrAuthPasswordInternalError
	}

	_, err = tx.Exec(`INSERT INTO users(id, name, role) 
	VALUES($1, $2, $3)`, user.ID, user.Name, user.Role)

	if err != nil {
		log.Error().Caller().Msg(err.Error())

		err = tx.Rollback()
		if err != nil {
			log.Error().Caller().Msg(err.Error())
		}

		return e.ErrAuthPasswordInternalError
	}

	_, err = tx.Exec(`INSERT INTO identities(id, user_id, provider, uid, password) 
		VALUES($1, $2, $3, $4, $5)`,
		identity.ID,
		user.ID,
		identity.Provider,
		identity.UID,
		identity.Password)

	if err != nil {
		log.Error().Caller().Msg(err.Error())

		err = tx.Rollback()
		if err != nil {
			log.Error().Caller().Msg(err.Error())
		}

		return e.ErrAuthPasswordInternalError
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return e.ErrAuthPasswordInternalError
	}

	return
}
