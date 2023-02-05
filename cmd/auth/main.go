package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/app/auth"
	"github.com/wvoliveira/corgi/internal/app/auth/facebook"
	"github.com/wvoliveira/corgi/internal/app/auth/google"
	"github.com/wvoliveira/corgi/internal/app/auth/password"
	"github.com/wvoliveira/corgi/internal/app/health"
	"github.com/wvoliveira/corgi/internal/pkg/constants"
	"github.com/wvoliveira/corgi/internal/pkg/database"
	"github.com/wvoliveira/corgi/internal/pkg/server"
)

func main() {
	db := database.NewSQL(flagDatasource)
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
		// Auth service: logout and check.
		service := auth.NewService(db)
		service.NewHTTP(apiRouter)
	}

	{
		// Auth password service.
		service := password.NewService(db, kv)
		service.NewHTTP(apiRouter)
	}

	{
		// Auth with Google provider.
		service := google.NewService(db)
		service.NewHTTP(apiRouter)
	}

	{
		// Auth with Facebook provider.
		service := facebook.NewService(db)
		service.NewHTTP(apiRouter)
	}

	{
		// Healthcheck endpoints.
		service := health.NewService(db, constants.VERSION)
		service.NewHTTP(rootRouter)
	}

	server.Graceful(router, flagHTTPPort)
}
