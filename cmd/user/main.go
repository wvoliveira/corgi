package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/app/health"
	"github.com/wvoliveira/corgi/internal/app/user"
	"github.com/wvoliveira/corgi/internal/pkg/constants"
	"github.com/wvoliveira/corgi/internal/pkg/database"
	"github.com/wvoliveira/corgi/internal/pkg/server"
)

func main() {
	db := database.NewSQL()
	kv := database.NewKV()

	// Create a root router and attach session.
	// I think its a good idea because we can manager user access with cookie based.
	router := gin.Default()

	server.AddStoreSession(router, flagSecretKey)

	rootRouter := router.Group("")
	apiRouter := router.Group("/api")

	if flagDebug {
		server.AddPProf(router, rootRouter)
	}

	{
		// User management service. Like profile view and edit.
		service := user.NewService(db, kv)
		service.NewHTTP(apiRouter)
	}

	{
		// Healthcheck endpoints.
		service := health.NewService(db, constants.VERSION)
		service.NewHTTP(rootRouter)
	}

	server.Graceful(router, flagHTTPPort)
}
