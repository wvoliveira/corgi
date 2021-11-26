package server

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/sethvargo/go-password/password"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type database struct {
	Logger zap.SugaredLogger
	DB     *gorm.DB
	Config Config
}

// NewDatabase create a gorm database object.
func NewDatabase(logger zap.SugaredLogger, config Config) database {
	return database{
		Logger: logger,
		DB:     initDatabase(logger, config),
		Config: config,
	}
}

func initDatabase(logger zap.SugaredLogger, config Config) (db *gorm.DB) {
	connString := fmt.Sprintf("postgres://%s@%s:%d/%s", config.DBUser, config.DBHost, config.DBPort, config.DBBase)
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		logger.Errorw("error configuring the database", "method", "initDatabase", "err", err.Error())
		os.Exit(1)
	}
	return db
}

// SeedUsers create the first users for system.
func (d database) SeedUsers() {
	accounts := []Account{
		{
			Name:  "Administrator",
			Email: "admin@local",
			Role:  "admin",
			Active: "true",
		},
		{
			Name:  "Normal user",
			Email: "user@local",
			Role:  "user",
			Active: "true",
		},
	}

	for _, acc := range accounts {
		if d.DB.Model(&acc).Where("email = ?", acc.Email).Updates(&acc).RowsAffected > 0 {
			return
		}

		secret, err := password.Generate(20, 5, 0, false, true)
		if err != nil {
			d.Logger.Errorw("fail to generate password", "method", "SeedUsers", "email", acc.Email, "err", err.Error())
			os.Exit(1)
		}

		messagePassword := fmt.Sprintf("account e-mail: %s password: %s", acc.Email, secret)
		d.Logger.Infow("password was generated", "method", "SeedUsers", "message", messagePassword)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(secret), 8)
		if err != nil {
			d.Logger.Errorw("error to generate hash password", "method", "SeeUsers", "err", err.Error())
			os.Exit(1)
		}

		acc.Password = string(hashedPassword)
		acc.ID = uuid.New().String()
		acc.CreatedAt = time.Now()
		d.DB.Model(&acc).Create(&acc)
	}
}
