package server

import (
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type cache struct {
	Logger zap.SugaredLogger
	DB     *redis.Client
	Config Config
}

// NewCache create a gorm cache object.
func NewCache(logger zap.SugaredLogger, config Config) cache {
	return cache{
		Logger: logger,
		DB:     initCache(logger, config),
		Config: config,
	}
}

func initCache(logger zap.SugaredLogger, config Config) (db *redis.Client) {
	return redis.NewClient(&redis.Options{
		Addr:     config.CacheAddr,
		Password: config.CachePassword,
		DB:       config.CacheDatabase,
	})
}
