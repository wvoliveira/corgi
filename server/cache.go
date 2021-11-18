package server

import (
	"github.com/go-kit/log"
	"github.com/go-redis/redis/v8"
)

type cache struct {
	Logger log.Logger
	DB     *redis.Client
	Config Config
}

// NewCache create a gorm cache object.
func NewCache(logger log.Logger, config Config) cache {
	return cache{
		Logger: logger,
		DB:     initCache(logger, config),
		Config: config,
	}
}

func initCache(logger log.Logger, config Config) (db *redis.Client) {
	return redis.NewClient(&redis.Options{
		Addr:     config.CacheAddr,
		Password: config.CachePassword,
		DB:       config.CacheDatabase,
	})
}
