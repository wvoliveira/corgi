package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/elga-io/redirect/app"

	"github.com/go-kit/kit/log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:embed ui/dist
//go:embed ui/dist/_next
//go:embed ui/dist/_next/static/chunks/pages/*.js
//go:embed ui/dist/_next/static/*/*.js
var nextFS embed.FS

func initialMigration(db *gorm.DB) {
	db.AutoMigrate(&app.URL{})
}

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

	// Init logging
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	// Init database
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		logger.Log("failed to connect database", err)
		os.Exit(2)
	}

	initialMigration(db)

	// Init services
	var s app.Service
	s = app.NewDBService(db)
	s = app.LoggingMiddleware(logger)(s)

	// Web UI
	distFS, err := fs.Sub(nextFS, "ui/dist")
	hh := http.FileServer(http.FS(distFS))

	if err != nil {
		logger.Log("exit", err)
		os.Exit(2)
	}

	var h http.Handler
	h = app.MakeHTTPHandler(hh, s, log.With(logger, "component", "HTTP"))

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
