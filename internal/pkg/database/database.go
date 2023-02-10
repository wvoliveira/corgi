// Package database provides DB transaction support for transactions tha span method calls of multiple
// repositories and services.
package database

import (
	"database/sql"
	"path/filepath"

	"github.com/dgraph-io/badger"
	"github.com/oklog/ulid/v2"
	"github.com/spf13/viper"
	"github.com/wvoliveira/corgi/internal/pkg/common"
	"github.com/wvoliveira/corgi/internal/pkg/model"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// NewSQL create a gorm database object.
func NewSQL() (db *sql.DB) {
	datasource := viper.GetString("datasource")
	db, err := sql.Open("postgres", datasource)
	if err != nil {
		panic("failed to connect in sqlite database")
	}
	return
}

// NewKV create a badger database object.
func NewKV() (db *badger.DB) {
	appFolder, err := common.GetOrCreateDataFolder()
	if err != nil {
		log.Fatal().Caller().Msg(err.Error())
	}

	dbFile := filepath.Join(appFolder, "cache")

	db, err = badger.Open(badger.DefaultOptions(dbFile))
	if err != nil {
		panic("failed to connect in sqlite database")
	}
	return
}

// SeedUsers create the first users for system.
func SeedUsers(db *sql.DB) {
	users := []model.User{}

	admin := model.User{
		ID:   ulid.Make().String(),
		Name: "Administrator",
		Role: "admin",
	}

	identityAdmin := model.Identity{
		ID:       ulid.Make().String(),
		UserID:   admin.ID,
		Provider: "email",
		UID:      "admin@corgi",
		Password: "password",
	}

	user := model.User{
		ID:   ulid.Make().String(),
		Name: "User",
		Role: "user",
	}

	identityUser := model.Identity{
		ID:       ulid.Make().String(),
		UserID:   user.ID,
		Provider: "email",
		UID:      "user@corgi",
		Password: "password",
	}

	admin.Identities = append(admin.Identities, identityAdmin)
	user.Identities = append(user.Identities, identityUser)
	users = append(users, user, admin)

	for _, u := range users {
		iden := u.Identities[0]

		// Check if provider and UID exists.
		var id string
		_ = db.QueryRow("SELECT id FROM identities WHERE provider = ? AND uid = ?").Scan(&id)
		if id != "" {
			continue
		}

		plainTextPassword := iden.Password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 8)
		iden.Password = string(hashedPassword)

		tx, err := db.Begin()
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return
		}

		_, err = tx.Exec(`INSERT INTO users(id, created_at, name, role) 
		VALUES(?, ?, ?, ?)`, u.ID, u.CreatedAt, u.Name, u.Role)

		if err != nil {
			log.Error().Caller().Msg(err.Error())

			err = tx.Rollback()
			if err != nil {
				log.Error().Caller().Msg(err.Error())
			}
			continue
		}

		_, err = tx.Exec(`INSERT INTO identities(id, user_id, created_at, provider, uid, password) 
		VALUES(?, ?, ?, ?, ?, ?)`,
			iden.ID,
			iden.UserID,
			iden.CreatedAt,
			iden.Provider,
			iden.UID,
			iden.Password,
		)

		if err != nil {
			log.Error().Caller().Msg(err.Error())

			err = tx.Rollback()
			if err != nil {
				log.Error().Caller().Msg(err.Error())
			}

			continue
		}

		err = tx.Commit()
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			continue
		}
	}
}
