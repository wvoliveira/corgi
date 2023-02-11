// Package database provides DB transaction support for transactions tha span method calls of multiple
// repositories and services.
package database

import (
	"database/sql"
	"path/filepath"

	"github.com/dgraph-io/badger"
	"github.com/spf13/viper"
	"github.com/wvoliveira/corgi/internal/pkg/common"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// NewSQL create a gorm database object.
func NewSQL() (db *sql.DB) {
	datasource := viper.GetString("datasource")
	db, err := sql.Open("postgres", datasource)
	if err != nil {
		panic("failed to connect in sqlite database")
	}
	return
}

// NewKV create a badger database object.
func NewKV() (db *badger.DB) {
	appFolder, err := common.GetOrCreateDataFolder()
	if err != nil {
		log.Fatal().Caller().Msg(err.Error())
	}

	dbFile := filepath.Join(appFolder, "cache")

	db, err = badger.Open(badger.DefaultOptions(dbFile))
	if err != nil {
		panic("failed to connect in sqlite database")
	}
	return
}
