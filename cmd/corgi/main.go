package main

import (
	"encoding/gob"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	cfa "github.com/naucon/casbin-fs-adapter"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/wvoliveira/corgi/configs/authorization"
	"github.com/wvoliveira/corgi/internal/app/auth/facebook"
	"github.com/wvoliveira/corgi/internal/app/auth/google"
	"github.com/wvoliveira/corgi/internal/app/auth/password"
	"github.com/wvoliveira/corgi/internal/app/auth/token"
	"github.com/wvoliveira/corgi/internal/app/click"
	"github.com/wvoliveira/corgi/internal/app/group"
	"github.com/wvoliveira/corgi/internal/app/health"
	"github.com/wvoliveira/corgi/internal/app/link"
	"github.com/wvoliveira/corgi/internal/app/user"
	"github.com/wvoliveira/corgi/internal/pkg/config"
	"github.com/wvoliveira/corgi/internal/pkg/constants"
	"github.com/wvoliveira/corgi/internal/pkg/database"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/middleware"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	ratelimit "github.com/wvoliveira/corgi/internal/pkg/rate-limit"
	"github.com/wvoliveira/corgi/internal/pkg/server"
	"github.com/wvoliveira/corgi/web"
	"net/http"
)

func init() {
	config.New()
	logger.Default()

	// Used as type for security session.
	// TODO: remove this because we use JWT now.
	gob.Register(model.User{})
}

func main() {
	db := database.NewSQL()
	cache := database.NewCache()

	// Create user Admin if not exists.
	// The password will be prompt to console at first run.
	database.CreateUserAdmin(db)

	// Enforce some authorization rules.
	// Check model and policy files in configs/authorization folder.
	authModel, _ := cfa.NewModel(authorization.DistFiles, "model.conf")
	authPolicy := cfa.NewAdapter(authorization.DistFiles, "policy.csv")
	enforcer, _ := casbin.NewEnforcer(authModel, authPolicy)

	// Create a root router and attach session.
	// I think it's a good idea because we can manage user access with cookie based.
	router := gin.New()
	router.Use(middleware.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.Authentication())
	router.Use(middleware.Authorization(enforcer))

	apiRouter := router.Group("/api")

	// Enable /metrics path Prometheus metrics like.
	// And middleware to add some basic metrics from routes.
	server.NewMetrics(apiRouter)

	// Don't enable some things with debug level.
	// Middleware for rate limit.
	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		server.AddPProf(apiRouter)
	} else {
		ratelimit.NewMiddleware(router, cache)
	}

	{
		// Auth service: logout and check.
		service := token.NewService(db)
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
		service := group.NewService(db, cache)
		service.NewHTTP(apiRouter)
	}

	{
		// Central business service: manage link shortener.
		service := link.NewService(db, cache)
		service.NewHTTP(router, apiRouter)
	}

	{
		// Clicks service. Metrics for each link.
		service := click.NewService(db, cache)
		service.NewHTTP(apiRouter)
	}

	{
		// Healthcheck endpoints.
		// Kubernetes healthcheck like: readiness and liveness.
		service := health.NewService(db, cache, constants.VERSION)
		service.NewHTTP(apiRouter)
	}

	// Send requests that do not have a router defined.
	router.NoRoute(func(c *gin.Context) {
		reqPath := c.Request.URL.Path
		c.FileFromFS(reqPath, http.FS(web.DistFS))
		return
	})

	server.Graceful(router, viper.GetInt("SERVER_HTTP_PORT"))
}
