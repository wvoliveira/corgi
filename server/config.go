package server

import (
	"os"

	"github.com/go-kit/log"
	"github.com/spf13/viper"
)

const (
	dbHost     = "localhost"
	dbPort     = 26257
	dbUser     = "root"
	dbPassword = ""
	dbBase     = "corgi"

	cacheAddr     = "localhost:6379"
	cachePassword = ""
	cacheDatabase = 0

	searchAddr     = "localhost:9200"
	searchPassword = ""

	serverAddress = "0.0.0.0:8000"
	secretKey     = "changeme"
)

// Config a struct for app configuration.
type Config struct {
	DBHost     string `mapstructure:"CORGI_DB_HOST"`
	DBPort     int    `mapstructure:"CORGI_DB_PORT"`
	DBUser     string `mapstructure:"CORGI_DB_USER"`
	DBPassword string `mapstructure:"CORGI_DB_PASSWORD"`
	DBBase     string `mapstructure:"CORGI_DB_BASE"`

	CacheAddr     string `mapstructure:"CORGI_CACHE_ADDR"`
	CachePassword string `mapstructure:"CORGI_CACHE_PASSWORD"`
	CacheDatabase int    `mapstructure:"CORGI_CACHE_DATABASE"`

	SearchAddr     string `mapstructure:"CORGI_SEARCH_ADDR"`
	SearchPassword string `mapstructure:"CORGI_SEARCH_PASSWORD"`

	ServerAddress string `mapstructure:"CORGI_SERVER_ADDRESS"`
	SecretKey     string `mapstructure:"CORGI_SECRET_KEY"`
}

// NewConfig load the configuration app.
func NewConfig(logger log.Logger, path string) (config Config) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	viper.ReadInConfig()

	// TODO: add rest of config.
	// viper "set default" doesnt work.
	if viper.Get("CORGI_DB_PORT") == 0 {
		viper.Set("CORGI_DB_PORT", dbPort)
	}
	if viper.Get("CORGI_DB_USER") == "" {
		viper.Set("CORGI_DB_USER", dbUser)
	}
	if viper.Get("CORGI_DB_BASE") == "" {
		viper.Set("CORGI_DB_BASE", dbBase)
	}
	if viper.Get("CORGI_SERVER_ADDRESS") == "" {
		viper.Set("CORGI_SERVER_ADDRESS", serverAddress)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		logger.Log("method", "NewConfig", "message", "error with viper.Unmarshal", "err", err.Error())
		os.Exit(1)
	}
	return
}
