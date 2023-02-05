package main

import (
	"encoding/gob"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

var (
	flagDebug      bool
	flagHTTPPort   string
	flagSecretKey  string
	flagDatasource string
)

func init() {
	flag.BoolVarP(&flagDebug, "debug", "d", false, "Enable DEBUG mode")
	flag.StringVar(&flagHTTPPort, "http-port", "8082", "Port for HTTP server")
	flag.StringVar(&flagSecretKey, "secret-key", "changeme", "Secret for encrypt session")
	flag.StringVar(&flagDatasource, "datasource", "postgres://user:password@localhost:5432/corgi?sslmode=disable", "URL for PostgreSQL")
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	if flagDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	gob.Register(model.User{})
}
