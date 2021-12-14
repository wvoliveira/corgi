package main

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/elga-io/corgi/internal/auth"
	"github.com/elga-io/corgi/internal/auth/facebook"
	"github.com/elga-io/corgi/internal/auth/google"
	"github.com/elga-io/corgi/internal/auth/password"
	"github.com/elga-io/corgi/internal/auth/token"
	"github.com/elga-io/corgi/internal/config"
	"github.com/elga-io/corgi/internal/entity"
	"github.com/elga-io/corgi/internal/health"
	"github.com/elga-io/corgi/internal/link"
	"github.com/elga-io/corgi/internal/public"
	"github.com/elga-io/corgi/internal/user"
	"github.com/elga-io/corgi/pkg/database"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// version indicates the current version of the application.
var version = "0.0.1"

func main() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Create root logger tagged with server version.
	logg := log.New().With(ctx, "version", version)

	// Load application configurations.
	cfg := config.NewConfig(logg, "configs")

	// Connect to the database and seed first users.
	db := database.NewDatabase(logg, cfg)
	if err := db.AutoMigrate(
		&entity.User{},
		&entity.Identity{},
		&entity.Link{},
		&entity.LinkLog{},
		&entity.Token{},
		&entity.Tag{},
		&entity.LocationIPv4{},
		&entity.LocationIPv6{},
	); err != nil {
		logg.Error("error in auto migrate", "err", err.Error())
		os.Exit(1)
	}
	database.SeedUsers(logg, db, cfg)

	// Setup Casbin auth rules.
	authEnforcer, err := casbin.NewEnforcer("./rbac_model.conf", "./rbac_policy.csv")
	if err != nil {
		logg.Error("error to get Casbin enforce rules")
		os.Exit(2)
	}

	// Start sessions.
	store := cookie.NewStore([]byte(cfg.App.SecretKey))

	// Auth services like login, register, logout, etc.
	authService := auth.NewService(logg, db, cfg.App.SecretKey, store, authEnforcer)
	authToken := token.NewService(logg, db, cfg.App.SecretKey, 30, store, authEnforcer)
	authPasswordService := password.NewService(logg, db, cfg.App.SecretKey, 30, store, authEnforcer)
	authGoogleService := google.NewService(logg, db, cfg, store, authEnforcer)
	authFacebookService := facebook.NewService(logg, db, cfg, store, authEnforcer)

	// Business services like links, users, etc.
	linkService := link.NewService(logg, db, cfg.App.SecretKey, store, authEnforcer)
	userService := user.NewService(logg, db, cfg.App.SecretKey, store, authEnforcer)

	// Public routes, like links?
	publicService := public.NewService(logg, db, store, authEnforcer)

	// Healthcheck services.
	healthService := health.NewService(logg, db, cfg.App.SecretKey, store, authEnforcer, version)

	// Initialize routers.
	router := gin.New()

	// Register business and needed routers.
	healthService.Routers(router)
	authService.Routers(router)
	authToken.Routers(router)
	authPasswordService.Routers(router)
	authGoogleService.Routers(router)
	authFacebookService.Routers(router)

	linkService.Routers(router)
	userService.Routers(router)
	publicService.Routers(router)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.HTTPPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		logg.Info("server listening :", cfg.Server.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Info("listen: %s", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	logg.Info("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logg.Info("Server forced to shutdown: ", err)
	}
	logg.Info("Server exiting")
}
