package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"

	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/wvoliveira/corgi/internal/app/auth/password"
	"github.com/wvoliveira/corgi/internal/app/debug"
	"github.com/wvoliveira/corgi/internal/pkg/database"
)

func main() {
	db := database.NewSQL()
	kv := database.NewKV()

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

	rootRouter := router.Group("")
	apiRouter := router.Group("/api")

	if flagDebug {
		pprof.Register(router)
		service := debug.NewService()
		service.NewHTTP(rootRouter)
	}

	{
		// Auth password service.
		service := password.NewService(db, kv)
		service.NewHTTP(apiRouter)
	}

	serverHTTPPort := viper.GetString("server.http_port")

	srv := &http.Server{
		Addr:         ":" + serverHTTPPort,
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
		log.Info().Caller().Msg(fmt.Sprintf("server listening http://127.0.0.1:%s", serverHTTPPort))

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
