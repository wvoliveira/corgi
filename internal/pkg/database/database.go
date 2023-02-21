// Package database provides DB transaction support for transactions tha span method calls of multiple
// repositories and services.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/wvoliveira/corgi/internal/pkg/common"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// NewSQL create a sql database object.
func NewSQL() (db *sql.DB) {
	log := logger.Logger(context.TODO())

	datasource := viper.GetString("DB_URL")
	db, err := sql.Open("postgres", datasource)

	if err != nil {
		log.Fatal().Caller().Msg(fmt.Sprintf("failed to connect to database: %s", err.Error()))
	}

	return
}

// NewKV create a cache/redis database object.
func NewCache() (db *redis.Client) {
	log := logger.Logger(context.TODO())

	datasource := viper.GetString("CACHE_URL")
	opt, err := redis.ParseURL(datasource)

	opt.DialTimeout = 3 * time.Second // no time unit = seconds
	opt.ReadTimeout = 6 * time.Second
	opt.MaxRetries = 2

	if err != nil {
		log.Fatal().Caller().Msg(fmt.Sprintf("failed to connect to cache: %s", err.Error()))
	}

	db = redis.NewClient(opt)
	status := db.Ping(context.TODO())

	if status.Err() != nil {
		log.Fatal().Caller().Msg(fmt.Sprintf("failed to connect to cache: %s", status.Err().Error()))
	}

	return
}

func CreateUserAdmin(db *sql.DB) {
	user := model.User{
		ID:       ulid.Make().String(),
		Name:     "Administrator",
		Username: "admin",
		Role:     "admin",
	}

	identity := model.Identity{
		ID:       ulid.Make().String(),
		UserID:   user.ID,
		Provider: "username",
		UID:      "admin",
	}

	// Check if provider and UID exists.
	var id string
	_ = db.QueryRow("SELECT id FROM identities WHERE provider = $1 AND uid = $2", identity.Provider, identity.UID).Scan(&id)
	if id != "" {
		return
	}

	plainTextPassword := common.CreateRandomPassword()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 8)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		os.Exit(2)
	}

	identity.Password = string(hashedPassword)

	tx, err := db.Begin()
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		os.Exit(2)
	}

	_, err = tx.Exec(`INSERT INTO users(id, name, username, role) VALUES($1, $2, $3, $4)`,
		user.ID, user.Name, user.Username, user.Role)

	if err != nil {
		log.Error().Caller().Msg("Error to create user admin: " + err.Error())

		err = tx.Rollback()
		if err != nil {
			log.Error().Caller().Msg(err.Error())
		}
		os.Exit(2)
	}

	_, err = tx.Exec(`INSERT INTO identities(id, user_id, provider, uid, password) 
	VALUES($1, $2, $3, $4, $5)`,
		identity.ID,
		identity.UserID,
		identity.Provider,
		identity.UID,
		identity.Password,
	)

	if err != nil {
		log.Error().Caller().Msg("Error to create user admin: " + err.Error())

		err = tx.Rollback()
		if err != nil {
			log.Error().Caller().Msg(err.Error())
		}
		os.Exit(2)
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	message := fmt.Sprintf(`
======================================
User "admin" created with successfull!
Username: admin
Password: %s
======================================
	`, plainTextPassword)

	log.Info().Caller().Msg(message)
}
