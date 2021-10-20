package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/elga-io/redir/urls"

	"github.com/go-kit/log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/patrickmn/go-cache"

	_ "github.com/elga-io/redir/urls/docs"
)

func initialMigration(db *gorm.DB) {
	db.AutoMigrate(&urls.URL{})
}

// @title URLs API
// @version 0.0.1
// @description Micro-serice for managing URLs
// @termsOfService http://elga.io/terms
// @contact.name API Support
// @contact.email support@elga.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
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

	// Init cache
	c := cache.New(5*time.Minute, 10*time.Minute)

	initialMigration(db)

	// Init services
	var s urls.Service
	s = urls.NewDBService(db, c)
	s = urls.LoggingMiddleware(logger)(s)

	if err != nil {
		logger.Log("exit", err)
		os.Exit(2)
	}

	var h http.Handler
	h = urls.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))

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
