package server

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/wvoliveira/corgi/internal/app/debug"
)

func Graceful(router *gin.Engine, httpPort int) {
	srv := &http.Server{
		Addr:         ":" + fmt.Sprintf("%d", httpPort),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Info().Caller().Msg(fmt.Sprintf("server listening http://127.0.0.1:%d", httpPort))

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

func AddStoreSession(router *gin.Engine) {
	secretKey := viper.GetString("secret_key")

	store := cookie.NewStore([]byte(secretKey))
	store.Options(sessions.Options{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60 * 24 * 30,
	})

	router.Use(sessions.Sessions("session", store))
}

func AddPProf(rg *gin.RouterGroup) {
	pprof.RouteRegister(rg)
	// Debug service like env vars, pprof route, etc.
	service := debug.NewService()
	service.NewHTTP(rg)
}

func NewMetrics(rg *gin.RouterGroup) {
	h := promhttp.Handler()
	hf := func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}

	rg.GET("/metrics", hf)
}
