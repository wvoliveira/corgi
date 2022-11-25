package cache

import (
	"time"

	"github.com/eko/gocache/v3/cache"
	"github.com/eko/gocache/v3/store"
	gocache "github.com/patrickmn/go-cache"
)

// New create a cache manager object with default values.
// Ex.: expires items in 5 minutes and purges/delete expired items in 10 minutes.
func New() *cache.Cache[[]byte] {
	gocacheClient := gocache.New(5*time.Minute, 10*time.Minute)
	gocacheStore := store.NewGoCache(gocacheClient)

	return cache.New[[]byte](gocacheStore)
}
