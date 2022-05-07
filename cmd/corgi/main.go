package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/elga-io/corgi/internal/app/auth"
	"github.com/elga-io/corgi/internal/app/auth/facebook"
	"github.com/elga-io/corgi/internal/app/auth/google"
	"github.com/elga-io/corgi/internal/app/auth/password"
	"github.com/elga-io/corgi/internal/app/auth/token"
	"github.com/elga-io/corgi/internal/app/config"
	"github.com/elga-io/corgi/internal/app/entity"
	"github.com/elga-io/corgi/internal/app/health"
	"github.com/elga-io/corgi/internal/app/link"
	"github.com/elga-io/corgi/internal/app/redirect"
	"github.com/elga-io/corgi/internal/app/user"
	"github.com/elga-io/corgi/internal/pkg/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// version indicates the current version of the application.
var version = "0.0.1"

//go:embed web
//go:embed web/_next/static
//go:embed web/_next/static/chunks/pages/*.js
//go:embed web/_next/static/*/*.js
var nextFS embed.FS

func main() {
	migrate := flag.Bool("m", false, "Enable GORM migration")
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Initializing app...")

	// Create database and cache folder in $HOME/.corgi path.
	folder, err := createDataFolder(".corgi")
	if err != nil {
		log.Fatal().Caller().Msg(err.Error())
	}

	var (
		cfg = config.NewConfig("configs")
		db  = database.NewSQLDatabase("sqlite", filepath.Join(folder, "data"))
	)

	if *migrate {
		// Seed first users. Most admins.
		if err := db.AutoMigrate(
			&entity.User{},
			&entity.Identity{},
			&entity.Link{},
			&entity.Token{},
		); err != nil {
			log.Fatal().Caller().Msg(err.Error())
			os.Exit(1)
		}
	}

	database.SeedUsers(db, cfg)

	// Setup Casbin auth rules.
	authEnforcer, err := casbin.NewEnforcer("./rbac_model.conf", "./rbac_policy.csv")
	if err != nil {
		log.Fatal().Caller().Msg("error to get Casbin enforce rules")
		os.Exit(2)
	}

	// Cors. Yes, we need this.
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:4200", "http://localhost:8081"}
	corsConfig.AddAllowMethods("*")
	corsConfig.AddAllowHeaders("*")
	corsConfig.AllowCredentials = true

	// mw := middleware.Middleware{Cache: cache, Logger: log.Logger}

	router := mux.NewRouter().SkipClean(true)
	// router.Use(mw.CorrelationID)
	// router.Use(mw.Log)

	apiRouter := router.PathPrefix("/api").Subrouter().StrictSlash(true)
	webRouter := router.PathPrefix("/").Subrouter().StrictSlash(true)

	// Start sessions.
	store := cookie.NewStore([]byte(cfg.App.SecretKey))

	{
		// Auth service: logout and check.
		service := auth.NewService(db, cfg.App.SecretKey, store, authEnforcer)
		service.NewHTTP(apiRouter)
	}

	{
		// Token refresh route.
		service := token.NewService(db, cfg.App.SecretKey, 30, store, authEnforcer)
		service.NewHTTP(apiRouter)
	}

	{
		// Auth password service.
		service := password.NewService(db, cfg.App.SecretKey, 30, store, authEnforcer)
		service.NewHTTP(apiRouter)
	}

	{
		// Auth with Google provider.
		service := google.NewService(db, cfg, store, authEnforcer)
		service.NewHTTP(apiRouter)
	}

	{
		// Auth with Facebook provider.
		service := facebook.NewService(db, cfg, store, authEnforcer)
		service.NewHTTP(apiRouter)
	}

	{
		// Central business service: manage link shortener.
		service := link.NewService(db, cfg.App.SecretKey, store, authEnforcer)
		service.NewHTTP(apiRouter)
	}

	{
		// User service. Like profile view and edit.
		service := user.NewService(db, cfg.App.SecretKey, store, authEnforcer)
		service.NewHTTP(apiRouter)
	}

	{
		// Central business service: redirect short link.
		service := redirect.NewService(db, store, authEnforcer)
		service.NewHTTP(apiRouter)
	}

	{
		// Healthcheck endpoints.
		service := health.NewService(db, cfg.App.SecretKey, store, authEnforcer, version)
		service.NewHTTP(apiRouter)
	}

	// Start web UI.
	distFS, err := fs.Sub(nextFS, "web")
	if err != nil {
		log.Fatal().Caller().Msg(err.Error())
		os.Exit(2)
	}
	webHandler := http.FileServer(http.FS(distFS))
	webRouter.PathPrefix("").Handler(webHandler)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.HTTPPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Info().Caller().Msg(fmt.Sprintf("server listening :%s", cfg.Server.HTTPPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Caller().Msg(err.Error())
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Info().Caller().Msg("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Info().Caller().Msg(fmt.Sprintf("Server forced to shutdown: %s", err.Error()))
	}
	log.Info().Caller().Msg("Server exiting")
}

func initWebUI() (ui http.FileSystem) {
	distFS, err := fs.Sub(nextFS, "web/dist")
	ui = http.FS(distFS)

	if err != nil {
		log.Info().Caller().Msg(err.Error())
		os.Exit(2)
	}
	return
}

func createDataFolder(name string) (folder string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	folder = filepath.Join(home, name)
	if _, err = os.Stat(folder); os.IsNotExist(err) {
		err = os.Mkdir(folder, os.ModePerm)
		if err != nil {
			return
		}
	} else if err != nil {
		return
	}
	return
}

func printRoutes(rs []*mux.Router) {
	for _, r := range rs {
		_ = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			uri, err := route.GetPathTemplate()
			if err != nil {
				log.Error().Msg(fmt.Sprintf("with get path template: %s", err.Error()))
				return err
			}

			method, err := route.GetMethods()
			if err != nil {
				if errors.Is(err, mux.ErrMethodMismatch) {
					return err
				}
			}

			if uri != "" && len(method) != 0 {
				log.Info().Msg(fmt.Sprintf("%s %s", uri, method))
			}
			return nil
		})
	}
}
