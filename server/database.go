package server

import (
	"os"

	"github.com/go-kit/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDatabase start database.
func InitDatabase(logger log.Logger, file *string) (db *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(*file), &gorm.Config{})
	if err != nil {
		logger.Log("failed to connect database", err)
		os.Exit(2)
	}
	return
}
