// Package database provides DB transaction support for transactions tha span method calls of multiple
// repositories and services.
package database

import (
	"os"
	"path/filepath"
	"time"

	"github.com/wvoliveira/corgi/internal/pkg/config"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	"github.com/wvoliveira/corgi/internal/pkg/util"
	"gorm.io/gorm/logger"

	"github.com/glebarez/sqlite"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// New create a gorm database object.
func New() (db *gorm.DB) {
	newLogger := logger.New(
		&log.Logger,
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
	cfg := gorm.Config{Logger: newLogger}

	// Create database and cache folder in $HOME/.corgi path.
	appFolder, err := util.CreateDataFolder(".corgi")
	if err != nil {
		log.Fatal().Caller().Msg(err.Error())
	}

	dbFile := filepath.Join(appFolder, "data")

	db, err = gorm.Open(sqlite.Open(dbFile), &cfg)
	if err != nil {
		panic("failed to connect in sqlite database")
	}
	return
}

// SeedUsers create the first users for system.
func SeedUsers(db *gorm.DB, c config.Config) {
	t := true
	users := []entity.User{
		{
			ID:        uuid.New().String(),
			CreatedAt: time.Now(),
			Name:      "Administrator",
			Role:      "admin",
			Active:    &t,
			Identities: []entity.Identity{
				{
					ID:        uuid.New().String(),
					CreatedAt: time.Now(),
					Provider:  "email",
					UID:       "admin@local",
					Password:  c.AdminPassword,
				},
			},
		},
		{
			ID:        uuid.New().String(),
			CreatedAt: time.Now(),
			Name:      "User",
			Role:      "user",
			Active:    &t,
			Identities: []entity.Identity{
				{
					ID:        uuid.New().String(),
					CreatedAt: time.Now(),
					Provider:  "email",
					UID:       "user@local",
					Password:  c.UserPassword,
				},
			},
		},
	}

	for _, user := range users {
		var count int64
		db.Model(&entity.Identity{}).Where("provider = ? AND uid = ?", user.Identities[0].Provider, user.Identities[0].UID).Count(&count)
		if count > 0 {
			continue
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Identities[0].Password), 8)
		if err != nil {
			log.Info().Caller().Msg(err.Error())
			os.Exit(1)
		}

		user.Identities[0].Password = string(hashedPassword)
		db.Model(&entity.User{}).Create(&user)
	}
}
