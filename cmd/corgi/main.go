package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wvoliveira/corgi/internal/app/auth"
	"github.com/wvoliveira/corgi/internal/app/auth/password"
	"github.com/wvoliveira/corgi/internal/app/jobs"
	"github.com/wvoliveira/corgi/internal/app/user"
	"github.com/wvoliveira/corgi/internal/pkg/cache"
	"github.com/wvoliveira/corgi/internal/pkg/config"
	"github.com/wvoliveira/corgi/internal/pkg/database"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
)

var (
	version       = "0.0.1"
	sessionSecret = "CHANGE_FOR_SOMETHING_MORE_SECURITY"

	flagDebug  bool
	flagConfig string
)

func init() {
	flag.BoolVar(&flagDebug, "debug", false, "Enable DEBUG mode")
	flag.StringVar(&flagConfig, "config", "~/.corgi/corgi.yaml", "Path of config file")
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if flagDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}

// //go:embed all:web
// var nextFS embed.FS

func main() {
	cfg := config.NewConfig(flagConfig)
	db := database.New()
	cache := cache.New()

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

	database.SeedUsers(db, cfg)

	// Create a root router and attach session.
	// I think its a good idea because we can manager user access with cookie based.
	router := gin.Default()

	store := cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60 * 24,
	})

	router.Use(sessions.Sessions("session", store))

	// rootRouter := router.Group("/")
	apiRouter := router.Group("/api")

	// rootRouter := router.PathPrefix("/").Subrouter().StrictSlash(true)
	// apiRouter := router.PathPrefix("/api").Subrouter().StrictSlash(true)
	// webRouter := router.MatcherFunc(func(req *http.Request, match *mux.RouteMatch) bool {
	// 	// Serve local web routes when either:
	// 	// - The request is for theses URIs:
	// 	// - / and /_next
	// 	return (req.URL.Path == "/" || strings.HasPrefix(req.URL.Path, "/_next"))
	// }).Subrouter().StrictSlash(true)

	if flagDebug {
		pprof.Register(router)
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
		// User service. Like profile view and edit.
		service := user.NewService(db, cache)
		service.NewHTTP(apiRouter)
	}

	// {
	// 	// Auth with Google provider.
	// 	service := google.NewService(db, cfg, store)
	// 	service.NewHTTP(apiRouter)
	// }

	// {
	// 	// Auth with Facebook provider.
	// 	service := facebook.NewService(db, cfg, store)
	// 	service.NewHTTP(apiRouter)
	// }

	// {
	// 	// Central business service: manage link shortener.
	// 	service := link.NewService(db, cfg.SecretKey, store)
	// 	service.NewHTTP(apiRouter)
	// }

	// {
	// 	// Healthcheck endpoints.
	// 	service := health.NewService(db, cfg, version)
	// 	service.NewHTTP(rootRouter)
	// }

	// {
	// 	// Info endpoint.
	// 	service := info.NewService(db, cfg, version)
	// 	service.NewHTTP(rootRouter)
	// }

	// {
	// 	// Central business service: redirect short link.
	// 	// Note: this service is on root router.
	// 	service := redirect.NewService(db, store)
	// 	service.NewHTTP(rootRouter)
	// }

	// {
	// 	// Start web application. User interface.
	// 	// Embedded UI.
	// 	distFS, err := fs.Sub(nextFS, "web")
	// 	if err != nil {
	// 		log.Fatal().Caller().Msg(err.Error())
	// 		os.Exit(2)
	// 	}

	// 	webHandler := http.FileServer(http.FS(distFS))
	// 	webRouter.PathPrefix("").Handler(webHandler)
	// }

	// Start cronjobs.
	serviceCron := jobs.NewService(db, cfg)
	serviceCron.Start()

	// Help func to get endpoints.
	// if flagDebug {
	// 	util.PrintRoutes([]*mux.Router{rootRouter, apiRouter})
	// }

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

	// Stop cronjobs.
	serviceCron.Stop()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Warn().Caller().Msg(fmt.Sprintf("server forced to shutdown: %s", err.Error()))
	}
	log.Info().Caller().Msg("server exited")
}
