package server

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// InitCache start cache layer.
func InitCache() (c *cache.Cache) {
	return cache.New(5*time.Minute, 10*time.Minute)

}
