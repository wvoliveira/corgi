package server

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Service store all methods. Yeah monolithic.
type Service struct {
	logger zap.SugaredLogger
	db     *gorm.DB
	secret string
}

// NewService create a new service with database and cache.
func NewService(logger zap.SugaredLogger, secretKey string, db *gorm.DB) Service {
	return Service{
		logger: logger,
		db:     db,
		secret: secretKey,
	}
}
