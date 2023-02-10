package password

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
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
	log := logger.Logger(c.Request.Context())
	idenDB := model.Identity{}

	query := "SELECT user_id, password FROM identities WHERE provider = $1 AND uid = $2"
	err = s.db.QueryRowContext(c, query, identity.Provider, identity.UID).Scan(&idenDB.UserID, &idenDB.Password)

	if err != nil {
		log.Warn().Caller().Msg(err.Error())
		return user, e.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(idenDB.Password), []byte(identity.Password)); err != nil {
		log.Info().Caller().Msg("authentication failed")
		return user, e.ErrUnauthorized
	}

	query = "SELECT * FROM users WHERE id = $1"
	err = s.db.QueryRowContext(c, query, idenDB.UserID).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Role, &user.Active)

	if err != nil {
		log.Warn().Caller().Msg(err.Error())
		return user, e.ErrAuthPasswordInternalError
	}

	return
}

// Register a new user to our database.
func (s service) Register(c *gin.Context, identity model.Identity) (err error) {
	log := logger.Logger(c.Request.Context())
	user := model.User{}
	identityDB := model.Identity{}

	query := "SELECT * FROM identities WHERE provider = $1 AND uid = $1"
	err = s.db.QueryRowContext(c, query, identity.Provider, identity.UID).Scan(&identityDB)

	if err == nil {
		log.Warn().Caller().Msg(fmt.Sprintf("provider '%s' and uid '%s' already exists", identity.Provider, identity.UID))
		return e.ErrAuthPasswordUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(identity.Password), 8)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return e.ErrAuthPasswordInternalError
	}

	identity.ID = ulid.Make().String()
	identity.CreatedAt = time.Now()
	identity.Password = string(hashedPassword)

	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.Role = "user"

	tx, err := s.db.BeginTx(c, &sql.TxOptions{})
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return e.ErrAuthPasswordInternalError
	}

	_, err = tx.Exec(`INSERT INTO identities(id, user_id, created_at, provider, uid, password) 
		VALUES(?, ?, ?, ?, ?, ?)`,
		identity.ID,
		user.ID,
		identity.CreatedAt,
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

	_, err = tx.Exec(`INSERT INTO users(id, created_at, name, role) 
	VALUES(?, ?, ?, ?)`, user.ID, user.CreatedAt, user.Name, user.Role)

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
