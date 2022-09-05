// Package database provides DB transaction support for transactions tha span method calls of multiple
// repositories and services.
package database

import (
	"os"
	"time"

	"github.com/wvoliveira/corgi/internal/app/config"
	"github.com/wvoliveira/corgi/internal/app/entity"
	"gorm.io/gorm/logger"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// NewSQLDatabase create a gorm database object.
// kind: type of database, like "sqlite", "mysql", "postgresql", etc.
// dsn: dsn with user/password and all necessary for connect in database.
func NewSQLDatabase(kind, dsn string) (db *gorm.DB) {
	newLogger := logger.New(
		&log.Logger, // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
	cfg := gorm.Config{Logger: newLogger}

	switch kind {
	case "sqlite":
		db, err := gorm.Open(sqlite.Open(dsn), &cfg)
		if err != nil {
			panic("failed to connect in sqlite database")
		}
		return db
	case "mysql":
		db, err := gorm.Open(mysql.Open(dsn), &cfg)
		if err != nil {
			panic("failed to connect in mysql database")
		}
		return db
	// case "postgresql":
	// 	db, err := gorm.Open(postgres.Open(dsn), &cfg)
	// 	if err != nil {
	// 		panic("failed to connect in postgresql database")
	// 	}
	// 	return db
	default:
		log.Fatal().Caller().Msg("this type of database is not supported")
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
					Password:  c.App.AdminPassword,
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
					Password:  c.App.UserPassword,
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
