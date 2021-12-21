package main

import (
	"context"
	"embed"
	"github.com/elga-io/corgi/pkg/log"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// version indicates the current version of the application.
var version = "0.0.1"

////go:embed web/dist
////go:embed web/dist/_next
////go:embed web/dist/_next/static/chunks/pages/*.js
////go:embed web/dist/_next/static/*/*.js
var nextFS embed.FS

func main() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Create root logger tagged with server version.
	logg := log.New().With(ctx, "version", version)

	ui := initWebUI(logg)
	//	middlewares := server.NewMiddlewares(*logger, config)

	mux := http.NewServeMux()
	mux.Handle("/", ui)

	//	http.Handle("/", middlewares.AccessControl(mux))
	srv := &http.Server{
		Addr:         ":8081",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		logg.Info("server listening :8081")
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

func initWebUI(logger log.Logger) (ui http.Handler) {
	distFS, err := fs.Sub(nextFS, "ui/dist")
	ui = http.FileServer(http.FS(distFS))

	if err != nil {
		logger.Info("error to start web UI", "err", err)
		os.Exit(2)
	}
	return
}
