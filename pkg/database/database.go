// Package database provides DB transaction support for transactions tha span method calls of multiple
// repositories and services.
package database

import (
	"fmt"
	"github.com/elga-io/corgi/internal/config"
	"github.com/elga-io/corgi/internal/entity"
	"github.com/elga-io/corgi/pkg/log"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDatabase create a gorm database object.
func NewDatabase(logger log.Logger, c config.Config) (db *gorm.DB) {
	return initDatabase(logger, c)
}

func initDatabase(logger log.Logger, c config.Config) (db *gorm.DB) {
	connString := fmt.Sprintf("postgres://%s@%s:%d/%s", c.Database.User, c.Database.Host, c.Database.Port, c.Database.Base)
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		logger.Error("error configuring the database", "method", "initDatabase", "err", err.Error())
		os.Exit(1)
	}
	return db
}

// SeedUsers create the first users for system.
func SeedUsers(logger log.Logger, db *gorm.DB, c config.Config) {
	users := []entity.User{
		{
			ID:        uuid.New().String(),
			CreatedAt: time.Now(),
			Name:      "Admin",
			Role:      "admin",
			Tags:      "admin,superuser,root",
			Active:    "true",
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
			Tags:      "user,limited",
			Active:    "true",
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
		db.Debug().Model(&entity.Identity{}).Where("provider = ? AND uid = ?", user.Identities[0].Provider, user.Identities[0].UID).Count(&count)
		if count > 0 {
			continue
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Identities[0].Password), 8)
		if err != nil {
			logger.Error("error to generate hash password", "method", "SeeUsers", "err", err.Error())
			os.Exit(1)
		}

		user.Identities[0].Password = string(hashedPassword)
		db.Debug().Model(&entity.User{}).Create(&user)
	}
}
