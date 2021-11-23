package server

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
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
func NewConfig(logger zap.SugaredLogger, path string) (config Config) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// Database config.
	if viper.Get("CORGI_DB_HOST") == nil {
		viper.Set("CORGI_DB_HOST", dbHost)
	}
	if viper.Get("CORGI_DB_PORT") == nil {
		viper.Set("CORGI_DB_PORT", dbPort)
	}
	if viper.Get("CORGI_DB_USER") == nil {
		viper.Set("CORGI_DB_USER", dbUser)
	}
	if viper.Get("CORGI_DB_PASSWORD") == nil {
		viper.Set("CORGI_DB_PASSWORD", dbPassword)
	}
	if viper.Get("CORGI_DB_BASE") == nil {
		viper.Set("CORGI_DB_BASE", dbBase)
	}

	// Cache config.
	if viper.Get("CORGI_CACHE_ADDR") == nil {
		viper.Set("CORGI_CACHE_ADDR", cacheAddr)
	}
	if viper.Get("CORGI_CACHE_PASSWORD") == nil {
		viper.Set("CORGI_CACHE_PASSWORD", cachePassword)
	}
	if viper.Get("CORGI_CACHE_DATABASE") == nil {
		viper.Set("CORGI_CACHE_DATABASE", cacheDatabase)
	}

	// Search engine config.
	if viper.Get("CORGI_SEARCH_ADDR") == nil {
		viper.Set("CORGI_SEARCH_ADDR", searchAddr)
	}
	if viper.Get("CORGI_SEARCH_PASSWORD") == nil {
		viper.Set("CORGI_SEARCH_PASSWORD", searchPassword)
	}

	// Service config.
	if viper.Get("CORGI_SERVER_ADDRESS") == nil {
		viper.Set("CORGI_SERVER_ADDRESS", serverAddress)
	}
	if viper.Get("CORGI_SECRET_KEY") == nil {
		viper.Set("CORGI_SECRET_KEY", secretKey)
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		logger.Warnf("error to read from config file", "method", "NewConfig", err.Error())
	}

	if err := viper.Unmarshal(&config); err != nil {
		logger.Errorf("error to load config", "method", "NewConfig", "err", err.Error())
		os.Exit(1)
	}
	return
}
