package server

import (
	"os"

	"github.com/go-kit/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDatabase start database.
func InitDatabase(logger log.Logger, config Config) (db *gorm.DB) {
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
