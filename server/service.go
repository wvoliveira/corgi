package server

import (
	"github.com/go-kit/log"
	"gorm.io/gorm"
)

// Service store all methods. Yeah monolithic.
type Service struct {
	logger log.Logger
	db     *gorm.DB
	secret string
}

// NewService create a new service with database and cache.
func NewService(logger log.Logger, secretKey string, db *gorm.DB) Service {
	return Service{
		logger: logger,
		db:     db,
		secret: secretKey,
	}
}
