package main

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/app/auth"
	"github.com/wvoliveira/corgi/internal/app/auth/facebook"
	"github.com/wvoliveira/corgi/internal/app/auth/google"
	"github.com/wvoliveira/corgi/internal/app/auth/password"
	"github.com/wvoliveira/corgi/internal/app/clicks"
	"github.com/wvoliveira/corgi/internal/app/group"
	"github.com/wvoliveira/corgi/internal/app/health"
	"github.com/wvoliveira/corgi/internal/app/link"
	"github.com/wvoliveira/corgi/internal/app/short"
	"github.com/wvoliveira/corgi/internal/app/user"
	"github.com/wvoliveira/corgi/internal/pkg/config"
	"github.com/wvoliveira/corgi/internal/pkg/constants"
	"github.com/wvoliveira/corgi/internal/pkg/database"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	ratelimit "github.com/wvoliveira/corgi/internal/pkg/rate-limit"
	"github.com/wvoliveira/corgi/internal/pkg/server"
	"github.com/wvoliveira/corgi/web"
)

func init() {
	config.New()
	logger.Default()
}

func main() {
	db := database.NewSQL()
	cache := database.NewCache()

	// Create user Admin if not exists.
	// You can desactive this user after installation!
	database.CreateUserAdmin(db)

	// Create a root router and attach session.
	// I think its a good idea because we can manager user access with cookie based.
	router := gin.Default()
	server.AddStoreSession(router)

	// First, check if request path is inside web app.
	// If yes, just answer the request and finish the request.
	router.Use(func(c *gin.Context) {
		reqPath := c.Request.URL.Path

		if !strings.HasPrefix(reqPath, "/api") {
			c.FileFromFS(reqPath, http.FS(web.DistFS))
			// c.Abort()
			return
		}
	})

	apiRouter := router.Group("/api")

	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		server.AddPProf(apiRouter)
	} else {
		// Dont enable some things with debug level.
		// Middleware for rate limit.
		ratelimit.NewMiddleware(router, cache)

	}

	{
		// Auth service: logout and check.
		service := auth.NewService(db)
		service.NewHTTP(apiRouter)
	}

	{
		// Auth password service.
		service := password.NewService(db, cache)
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
		// User management service. Like profile view and edit.
		service := user.NewService(db, cache)
		service.NewHTTP(apiRouter)
	}

	{
		// Group management service. Like create group, add users, etc.
		service := group.NewService(db)
		service.NewHTTP(apiRouter)
	}

	{
		// Central business service: manage link shortener.
		service := link.NewService(db, cache)
		service.NewHTTP(apiRouter)
	}

	{
		// Clicks service. Metrics for each link.
		service := clicks.NewService(db, cache)
		service.NewHTTP(apiRouter)
	}

	{
		// Healthcheck endpoints.
		service := health.NewService(db, constants.VERSION)
		service.NewHTTP(apiRouter)
	}

	{
		// Central business service: redirect short link.
		// Note: this service is on root router.
		service := short.NewService(db, cache)
		service.NewHTTP(apiRouter)
	}

	server.Graceful(router, viper.GetInt("SERVER_HTTP_PORT"))
}
