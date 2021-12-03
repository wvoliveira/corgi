package main

import (
	"context"
	"github.com/elga-io/corgi/internal/auth/password"
	"github.com/elga-io/corgi/internal/config"
	"github.com/elga-io/corgi/internal/entity"
	"github.com/elga-io/corgi/internal/health"
	"github.com/elga-io/corgi/internal/link"
	"github.com/elga-io/corgi/pkg/database"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/elga-io/corgi/pkg/middlewares"
	"github.com/gin-contrib/sessions"
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
	if err := db.AutoMigrate(&entity.User{}, &entity.Identity{}, &entity.Link{}, &entity.Token{}); err != nil {
		logg.Error("error in auto migrate", "err", err.Error())
		os.Exit(1)
	}
	database.SeedUsers(logg, db, cfg)

	// Services.
	authPasswordService := password.NewService(logg, db, cfg.App.SecretKey, 30)
	linkService := link.NewService(logg, db)
	healthService := health.NewService(logg, db, version)

	// Start sessions.
	store := cookie.NewStore([]byte(cfg.App.SecretKey))

	// Routers.
	router := gin.New()
	root := router.Group("/")

	// Handlers in root path ("/").
	healthService.Routers(root)

	api := router.Group("/api")
	api.Use(sessions.SessionsMany([]string{"session_unique", "session_auth"}, store))
	api.Use(middlewares.Access(logg))
	api.Use(middlewares.Checks())

	// Handlers in /api path.
	authPasswordService.Routers(api)

	v1 := router.Group("/api/v1")
	v1.Use(sessions.SessionsMany([]string{"session_unique", "session_auth"}, store))
	v1.Use(middlewares.Auth(logg, cfg.App.SecretKey))
	api.Use(middlewares.Access(logg))
	api.Use(middlewares.Checks())

	// Handlers in /api/v1 path.
	linkService.Routers(v1)

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
