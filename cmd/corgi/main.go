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
	"github.com/wvoliveira/corgi/internal/app/health"
	"github.com/wvoliveira/corgi/internal/app/link"
	appLog "github.com/wvoliveira/corgi/internal/app/log"
	"github.com/wvoliveira/corgi/internal/app/redirect"
	"github.com/wvoliveira/corgi/internal/app/user"
	"github.com/wvoliveira/corgi/internal/pkg/constants"
	"github.com/wvoliveira/corgi/internal/pkg/database"
	"github.com/wvoliveira/corgi/internal/pkg/server"
	"github.com/wvoliveira/corgi/web"
)

func main() {
	db := database.NewSQL(flagDatasource)
	kv := database.NewKV()

	// Seed first users. Most admins.
	database.SeedUsers(db)

	// Create a root router and attach session.
	// I think its a good idea because we can manager user access with cookie based.
	router := gin.Default()

	server.AddStoreSession(router, viper.GetString("secret_key"))

	// First, check if request path is inside web app.
	// If yes, just answer the request and finish the request.
	router.Use(func(c *gin.Context) {
		reqPath := c.Request.URL.Path

		if reqPath == "/" {
			c.FileFromFS(reqPath, http.FS(web.DistFS))
			c.Abort()
		}

		webPrefixPaths := []string{
			"/_next", "/favicon.ico", "/search", "/login", "/register", "/settings", "/profile",
		}

		for _, path := range webPrefixPaths {

			if strings.HasPrefix(reqPath, path) {
				router.RedirectTrailingSlash = false

				c.FileFromFS(reqPath, http.FS(web.DistFS))
				c.Abort()
				return
			}

		}
	})

	apiRouter := router.Group("/api")

	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		server.AddPProf(router, apiRouter)
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
		// User management service. Like profile view and edit.
		service := user.NewService(db, kv)
		service.NewHTTP(apiRouter)
	}

	{
		// Central business service: manage link shortener.
		service := link.NewService(db)
		service.NewHTTP(apiRouter)
	}

	{
		// Clicks service. Metrics for each link.
		service := clicks.NewService(db, kv)
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
		service := redirect.NewService(db, kv)
		service.NewHTTP(apiRouter)
	}

	{
		// Central business service: redirect short link.
		// Note: this service is on root router.
		service := appLog.NewService()
		service.NewHTTP(apiRouter)
	}

	server.Graceful(router, viper.GetString("server.http_port"))
}
