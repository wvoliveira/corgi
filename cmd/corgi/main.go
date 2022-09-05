package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wvoliveira/corgi/internal/app/auth"
	"github.com/wvoliveira/corgi/internal/app/auth/facebook"
	"github.com/wvoliveira/corgi/internal/app/auth/google"
	"github.com/wvoliveira/corgi/internal/app/auth/password"
	"github.com/wvoliveira/corgi/internal/app/auth/token"
	"github.com/wvoliveira/corgi/internal/app/config"
	"github.com/wvoliveira/corgi/internal/app/entity"
	"github.com/wvoliveira/corgi/internal/app/health"
	"github.com/wvoliveira/corgi/internal/app/link"
	"github.com/wvoliveira/corgi/internal/app/redirect"
	"github.com/wvoliveira/corgi/internal/app/user"
	"github.com/wvoliveira/corgi/internal/pkg/database"
	"github.com/wvoliveira/corgi/internal/pkg/middleware"
	"github.com/wvoliveira/corgi/internal/pkg/util"
)

// version indicates the current version of the application.
var version = "0.0.1"

//go:embed all:web
var nextFS embed.FS

func main() {
	debug := flag.Bool("d", false, "Enable DEBUG mode")
	migrate := flag.Bool("m", true, "Enable GORM migration")
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Create database and cache folder in $HOME/.corgi path.
	folder, err := util.CreateDataFolder(".corgi")
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

	// TODO: setup a properly CORS.

	router := mux.NewRouter().SkipClean(false)
	router.Use(middleware.Access)

	rootRouter := router.PathPrefix("/").Subrouter().StrictSlash(true)
	apiRouter := router.PathPrefix("/api").Subrouter().StrictSlash(true)
	webRouter := router.MatcherFunc(func(req *http.Request, match *mux.RouteMatch) bool {
		// Serve local web routes when either:
		// - The request is for theses URIs:
		// - / and /_next
		return (req.URL.Path == "/" || strings.HasPrefix(req.URL.Path, "/_next"))
	}).Subrouter().StrictSlash(true)

	// Start sessions.
	store := sessions.NewCookieStore([]byte(cfg.App.SecretKey))

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
		// Healthcheck endpoints.
		service := health.NewService(db, cfg.App.SecretKey, store, authEnforcer, version)
		service.NewHTTP(rootRouter)
	}

	{
		// Central business service: redirect short link.
		// Note: this service is on root router.
		service := redirect.NewService(db, store, authEnforcer)
		service.NewHTTP(rootRouter)
	}

	// Start web UI.
	distFS, err := fs.Sub(nextFS, "web")
	if err != nil {
		log.Fatal().Caller().Msg(err.Error())
		os.Exit(2)
	}

	// Web UI embedded.
	webHandler := http.FileServer(http.FS(distFS))
	webRouter.PathPrefix("").Handler(webHandler)

	// Help func to get endpoints.
	if *debug {
		util.PrintRoutes([]*mux.Router{rootRouter, apiRouter})
	}

	srv := &http.Server{
		Addr:         ":" + cfg.Server.HTTPPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Info().Caller().Msg(fmt.Sprintf("server listening http://127.0.0.1:%s", cfg.Server.HTTPPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Caller().Msg(err.Error())
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Info().Caller().Msg("shutting down gracefully..")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Warn().Caller().Msg(fmt.Sprintf("server forced to shutdown: %s", err.Error()))
	}
	log.Info().Caller().Msg("server exited")
}
