package ratelimit

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"

	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

func NewMiddleware(router *gin.Engine, cache *redis.Client) {
	// Define a limit rate to 10 requests per seconds.
	rate, err := limiter.NewRateFromFormatted("10-S")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create a store with the redis client.
	store, err := sredis.NewStoreWithOptions(cache, limiter.StoreOptions{
		Prefix:   "rate_limit",
		MaxRetry: 3,
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	// Create a new middleware with the limiter instance.
	middleware := mgin.NewMiddleware(limiter.New(store, rate))

	router.ForwardedByClientIP = true
	router.Use(middleware)
}
