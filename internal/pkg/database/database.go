// Package database provides DB transaction support for transactions tha span method calls of multiple
// repositories and services.
package database

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/wvoliveira/corgi/internal/pkg/common"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	"gorm.io/gorm"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// NewSQL create a gorm database object.
func NewSQL(datasource string) (db *sql.DB) {
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
func SeedUsers(db *gorm.DB) {
	t := true
	users := []model.User{
		{
			ID:        uuid.New().String(),
			CreatedAt: time.Now(),
			Name:      "Administrator",
			Role:      "admin",
			Active:    &t,
			Identities: []model.Identity{
				{
					ID:        uuid.New().String(),
					CreatedAt: time.Now(),
					Provider:  "email",
					UID:       "admin@corgi",
					Password:  "password",
				},
			},
		},
		{
			ID:        uuid.New().String(),
			CreatedAt: time.Now(),
			Name:      "User",
			Role:      "user",
			Active:    &t,
			Identities: []model.Identity{
				{
					ID:        uuid.New().String(),
					CreatedAt: time.Now(),
					Provider:  "email",
					UID:       "user@corgi",
					Password:  "password",
				},
			},
		},
	}

	for _, user := range users {
		var count int64

		provider := user.Identities[0].Provider
		uid := user.Identities[0].UID

		db.Model(&model.Identity{}).
			Where("provider = ? AND uid = ?", provider, uid).
			Count(&count)

		if count > 0 {
			continue
		}

		plainTextPassword := user.Identities[0].Password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 8)

		if err != nil {
			log.Info().Caller().Msg(err.Error())
			os.Exit(1)
		}

		user.Identities[0].Password = string(hashedPassword)
		db.Model(&model.User{}).Create(&user)
	}
}
