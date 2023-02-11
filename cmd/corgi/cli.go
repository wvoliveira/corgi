package main

import (
	"encoding/gob"
	"fmt"

	flag "github.com/spf13/pflag"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

var (
	flagLogLevel    string
	flagRedirectURL string
	flagDatasource  string

	flagServerHTTPPort     int
	flagServerReadTimeout  int
	flagServerWriteTimeout int

	flagDomainDefault      string
	flagDomainAlternatives []string
)

func init() {
	flag.StringVar(&flagLogLevel, "log-level", "info", "Log level")
	flag.StringVar(&flagRedirectURL, "redirect-url", "http://127.0.0.1:8081", "Redirect URL for OAUTH")
	flag.StringVar(&flagDatasource, "datasource", "postgres://user:password@localhost:5432/corgi?sslmode=disable", "URL for PostgreSQL")

	flag.IntVar(&flagServerHTTPPort, "server-http-port", 8081, "")
	flag.IntVar(&flagServerReadTimeout, "server-read-timeout", 10, "")
	flag.IntVar(&flagServerWriteTimeout, "server-write-timeout", 10, "")

	flag.StringVar(&flagDomainDefault, "domain-default", fmt.Sprintf("localhost:%d", flagServerHTTPPort), "")
	flag.StringArrayVar(&flagDomainAlternatives, "domain-alternatives", []string{}, "")
	flag.Parse()

	logger.Default()

	gob.Register(model.User{})
}
