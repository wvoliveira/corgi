package main

import (
	"embed"
	"fmt"
	"github.com/elga-io/corgi/pkg/log"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// version indicates the current version of the application.
var version = "0.0.1"

//go:embed web/dist
//go:embed web/dist/_next
//go:embed web/dist/_next/static/chunks/pages/*.js
//go:embed web/dist/_next/static/*/*.js
var nextFS embed.FS

func main() {
	//	logger := server.NewLogger()
	logg := log.New().With(nil, "version", version)
	//	config := server.NewConfig(*logger, ".")

	ui := initWebUI(logg)
	//	middlewares := server.NewMiddlewares(*logger, config)

	mux := http.NewServeMux()
	mux.Handle("/", ui)

	//	http.Handle("/", middlewares.AccessControl(mux))
	startServer(logg, ":8080")
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

func startServer(logger log.Logger, httpAddr string) {
	errs := make(chan error, 2)
	go func() {
		logger.Info("listening http server", "address", httpAddr)
		errs <- http.ListenAndServe(httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Info("terminated", "err", <-errs)
}
