// Package database provides DB transaction support for transactions tha span method calls of multiple
// repositories and services.
package database

import (
	"database/sql"

	"github.com/spf13/viper"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// NewSQL create a sql database object.
func NewSQL() (db *sql.DB) {
	datasource := viper.GetString("DB_URL")
	db, err := sql.Open("postgres", datasource)

	if err != nil {
		panic("failed to connect in sqlite database")
	}

	return
}

// NewKV create a cache/redis database object.
func NewCache() (db *redis.Client) {
	datasource := viper.GetString("CACHE_URL")
	opt, err := redis.ParseURL(datasource)

	if err != nil {
		panic(err)
	}

	return redis.NewClient(opt)
}
