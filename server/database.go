package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

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

func (d database) getAccountKey(id, email string) (dbKey, cacheKey string) {
	return fmt.Sprintf("db_account_id:%s:account_email:%s", id, email),
		fmt.Sprintf("cache_account_id:%s:account_email:%s", id, email)
}

// Create the first user for system.
func (d database) SeedUsers() {
	var (
		account Account
		ctx     = context.Background()
	)

	account = Account{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		LastLogin: time.Now(),

		Name:  "Administrator",
		Email: "admin@local",
		Role:  "admin",

		Active: "true",
	}

	dbKeyPattern, _ := d.getAccountKey("*", account.Email)

	keys, _ := d.DB.Keys(ctx, dbKeyPattern).Result()
	if len(keys) != 0 {
		return
	}

	// Generate a random password for admin account.
	secret, err := password.Generate(20, 5, 0, false, true)
	if err != nil {
		d.Logger.Log("method", "SeedUsers", "message", "fail to generate password", "err", err.Error())
		os.Exit(1)
	}

	messagePassword := fmt.Sprintf("admin user: admin@local password: %s", secret)
	d.Logger.Log("method", "SeedUsers", "message", messagePassword)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(secret), 8)
	if err != nil {
		d.Logger.Log("method", "SeedUsers", "message", "fail to hash password", "err", err.Error())
		os.Exit(1)
	}

	account.Password = string(hashedPassword)

	accountJs, err := json.Marshal(account)
	if err != nil {
		d.Logger.Log("method", "SeedUsers", "message", "error to marshal account to json", "err", err.Error())
	}

	key, _ := d.getAccountKey(account.ID, account.Email)

	if err = d.DB.Set(ctx, key, accountJs, 0).Err(); err != nil {
		d.Logger.Log("method", "SeedUsers", "message", "fail to create admin account", "err", err.Error())
		os.Exit(1)
	}
}
