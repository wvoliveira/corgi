package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
)

type database struct {
	Logger log.Logger
	DB     *redis.Client
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

func initDatabase(logger log.Logger, config Config) (db *redis.Client) {
	return redis.NewClient(&redis.Options{
		Addr:     config.DBAddr,
		Password: config.DBPassword,
		DB:       config.DBDatabase,
	})
}

// Create the first user for system.
func (d database) SeedUsers() {
	var account Account
	var ctx = context.Background()

	// Check if admin account exists.
	key := fmt.Sprintf("db_account_email:%s", "admin@local")
	if _, err := d.DB.Get(ctx, key).Result(); err != redis.Nil {
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
		ID:       uuid.New().String(),
		Name:     "admin",
		Email:    "admin@local",
		Password: string(hashedPassword),
	}

	accountJs, err := json.Marshal(account)
	if err != nil {
		d.Logger.Log("method", "SeedUsers", "message", "error to marshal account to json", "err", err.Error())
	}

	if err = d.DB.Set(ctx, key, accountJs, 0).Err(); err != nil {
		d.Logger.Log("method", "SeedUsers", "message", "fail to create admin account", "err", err.Error())
		os.Exit(1)
	}
}
