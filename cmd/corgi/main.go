package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wvoliveira/corgi/internal/app/auth"
	"github.com/wvoliveira/corgi/internal/app/auth/facebook"
	"github.com/wvoliveira/corgi/internal/app/auth/google"
	"github.com/wvoliveira/corgi/internal/app/auth/password"
	"github.com/wvoliveira/corgi/internal/app/clicks"
	"github.com/wvoliveira/corgi/internal/app/debug"
	"github.com/wvoliveira/corgi/internal/app/health"
	"github.com/wvoliveira/corgi/internal/app/link"
	"github.com/wvoliveira/corgi/internal/app/redirect"
	"github.com/wvoliveira/corgi/internal/app/user"
	"github.com/wvoliveira/corgi/internal/pkg/config"
	"github.com/wvoliveira/corgi/internal/pkg/database"
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

var (
	version = "0.0.1"

	flagDebug  bool
	flagConfig string
)

func init() {
	flag.BoolVarP(&flagDebug, "debug", "d", false, "Enable DEBUG mode")
	flag.StringVarP(&flagConfig, "config", "c", "~/.corgi/corgi.yaml", "Path of config file")
	flag.Parse()

	config.New(flagConfig)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if flagDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	gob.Register(model.User{})
}

// //go:embed all:web
// var nextFS embed.FS

func main() {
	db := database.NewSQL()
	kv := database.NewKV()

	// Seed first users. Most admins.
	if err := db.AutoMigrate(
		&model.User{},
		&model.Identity{},
		&model.Link{},
		&model.Token{},
	); err != nil {
		log.Fatal().Caller().Msg(err.Error())
		os.Exit(1)
	}

	database.SeedUsers(db)

	// Create a root router and attach session.
	// I think its a good idea because we can manager user access with cookie based.
	router := gin.Default()

	store := cookie.NewStore([]byte(viper.GetString("secret_key")))
	store.Options(sessions.Options{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60 * 24 * 30,
	})

	router.Use(sessions.Sessions("session", store))

	rootRouter := router.Group("/")
	apiRouter := rootRouter.Group("/api")

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
		service := debug.NewService()
		service.NewHTTP(rootRouter)
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
		// User management service. Like profile view and edit.
		service := user.NewService(db, kv)
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
		service := health.NewService(db, version)
		service.NewHTTP(rootRouter)
	}

	{
		// Central business service: redirect short link.
		// Note: this service is on root router.
		service := redirect.NewService(db, kv)
		service.NewHTTP(rootRouter)
	}

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

	srv := &http.Server{
		Addr:         ":" + viper.GetString("server.http_port"),
		Handler:      router,
		ReadTimeout:  viper.GetDuration("server.read_timeout") * time.Second,
		WriteTimeout: viper.GetDuration("server.write_timeout") * time.Second,
	}

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Info().Caller().Msg(fmt.Sprintf("server listening http://127.0.0.1:%s", viper.GetString("server.http_port")))

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
