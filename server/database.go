package server

import (
	"fmt"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/google/uuid"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
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
	connstring := fmt.Sprintf("postgres://%s@%s:%d/%s", config.DBUser, config.DBHost, config.DBPort, config.DBBase)
	db, err := gorm.Open(postgres.Open(connstring), &gorm.Config{})
	if err != nil {
		logger.Log("method", "initDatabase", "message", "error configuring the database", "err", err.Error())
		os.Exit(0)
	}
	return db
}

// Create the first users for system.
func (d database) SeedUsers() {
	d.addAccountAdmin()
	d.addAccountUser()
}

func (d database) addAccountAdmin() {
	account := Account{
		Name:  "Administrator",
		Email: "admin@local",
		Role:  "admin",

		Active: "true",
	}

	if d.DB.Model(&account).Where("email = ?", account.Email).Updates(&account).RowsAffected > 0 {
		return
	}

	secret, err := password.Generate(20, 5, 0, false, true)
	if err != nil {
		d.Logger.Log("method", "addAccountAdmin", "message", "fail to generate password", "err", err.Error())
		os.Exit(1)
	}

	messagePassword := fmt.Sprintf("account admin: admin@local password: %s", secret)
	d.Logger.Log("method", "addAccountAdmin", "message", messagePassword)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(secret), 8)
	if err != nil {
		d.Logger.Log("method", "addAccountAdmin", "message", "fail to hash password", "err", err.Error())
		os.Exit(1)
	}

	account.Password = string(hashedPassword)

	account.ID = uuid.New().String()
	account.CreatedAt = time.Now()

	d.DB.Model(&account).Create(&account)
}

func (d database) addAccountUser() {
	account := Account{
		Name:  "Normal user",
		Email: "user@local",
		Role:  "user",

		Active: "true",
	}

	if d.DB.Model(&account).Where("email = ?", account.Email).Updates(&account).RowsAffected > 0 {
		return
	}

	secret, err := password.Generate(20, 5, 0, false, true)
	if err != nil {
		d.Logger.Log("method", "addAccountUser", "message", "fail to generate password", "err", err.Error())
		os.Exit(1)
	}

	messagePassword := fmt.Sprintf("account user: user@local password: %s", secret)
	d.Logger.Log("method", "addAccountUser", "message", messagePassword)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(secret), 8)
	if err != nil {
		d.Logger.Log("method", "addAccountUser", "message", "fail to hash password", "err", err.Error())
		os.Exit(1)
	}

	account.Password = string(hashedPassword)

	account.ID = uuid.New().String()
	account.CreatedAt = time.Now()

	d.DB.Model(&account).Create(&account)
}
