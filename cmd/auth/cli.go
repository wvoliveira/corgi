package main

import (
	"encoding/gob"

	flag "github.com/spf13/pflag"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
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
	flag.StringVar(&flagHTTPPort, "http-port", "8081", "Port for HTTP server")
	flag.StringVar(&flagSecretKey, "secret-key", "changeme", "Secret for encrypt session")
	flag.StringVar(&flagDatasource, "datasource", "postgres://user:password@localhost:5432/corgi?sslmode=disable", "URL for PostgreSQL")
	flag.Parse()

	logger.Default(flagDebug)

	gob.Register(model.User{})
}
