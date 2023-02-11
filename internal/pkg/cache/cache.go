package cache

import (
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/dgraph-io/badger"
	"github.com/wvoliveira/corgi/internal/pkg/common"
)

// New create a cache manager object with default values.
// Ex.: expires items in 5 minutes and purges/delete expired items in 10 minutes.
func New() (db *badger.DB) {
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
