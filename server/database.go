package server

import (
	"fmt"
	"os"

	"github.com/go-kit/log"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type database struct {
	Logger log.Logger
	DB     *gorm.DB
	Config Config
}

// NewDatabase create a gorm database object.
func NewDatabase(logger log.Logger, config Config) database {
	return database{
		Logger: logger,
		DB:     initDatabase(logger, config),
		Config: config,
	}
}

func initDatabase(logger log.Logger, config Config) (db *gorm.DB) {
	var err error

	switch dbType := config.DBType; dbType {
	case "sqlite":
		db, err = loadSqlite(config)
	default:
		db, err = loadSqlite(config)
	}

	if err != nil {
		logger.Log("failed to connect database", err)
		os.Exit(2)
	}
	return
}

func loadSqlite(config Config) (db *gorm.DB, err error) {
	if config.DBType == "memory" {
		db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		return
	}

	db, err = gorm.Open(sqlite.Open(config.DBSource), &gorm.Config{})
	return
}

// Create the first user for system.
func (d database) SeedUsers() {
	var account Account

	// Check if admin account exists.
	if err := d.DB.Where("Email = ?", "admin@local").First(&account).Error; err == nil {
		return
	}

	// Generate a random password for admin account.
	secret, err := password.Generate(20, 5, 0, false, true)
	if err != nil {
		d.Logger.Log("method", "SeedUsers", "message", "fail to generate password", "err", err.Error())
		os.Exit(1)
	}

	messagePassword := fmt.Sprintf("admin password: %s", secret)
	d.Logger.Log("method", "SeedUsers", "message", messagePassword)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(secret), 8)
	if err != nil {
		d.Logger.Log("method", "SeedUsers", "message", "fail to hash password", "err", err.Error())
		os.Exit(1)
	}

	account = Account{
		Name:     "admin",
		Email:    "admin@local",
		Password: string(hashedPassword),
	}

	if err = d.DB.Model(account).Create(&account).Error; err != nil {
		d.Logger.Log("method", "SeedUsers", "message", "fail to create admin account", "err", err.Error())
		os.Exit(1)
	}
}
