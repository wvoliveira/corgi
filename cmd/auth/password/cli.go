package main

import (
	"encoding/gob"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"
	"github.com/wvoliveira/corgi/internal/pkg/config"
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

var (
	flagDebug  bool
	flagConfig string
)

func init() {
	flag.BoolVarP(&flagDebug, "debug", "d", false, "Enable DEBUG mode")
	flag.StringVarP(&flagConfig, "config", "c", "~/.corgi/corgi_auth_password.yaml", "Path of config file")
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
