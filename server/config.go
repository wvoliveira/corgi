package server

import (
	"os"

	"github.com/go-kit/log"
	"github.com/spf13/viper"
)

const (
	dbAddr     = "localhost:6379"
	dbPassword = ""
	dbDatabase = 0

	cacheAddr     = "localhost:6379"
	cachePassword = ""
	cacheDatabase = 1

	serverAddress = "0.0.0.0:8080"
	secretKey     = "changeme"
)

// Config a struct for app configuration.
type Config struct {
	DBAddr     string `mapstructure:"REDIR_DB_ADDR"`
	DBPassword string `mapstructure:"REDIR_DB_PASSWORD"`
	DBDatabase int    `mapstructure:"REDIR_DB_DATABASE"`

	CacheAddr     string `mapstructure:"REDIR_CACHE_ADDR"`
	CachePassword string `mapstructure:"REDIR_CACHE_PASSWORD"`
	CacheDatabase int    `mapstructure:"REDIR_CACHE_DATABASE"`

	ServerAddress string `mapstructure:"REDIR_SERVER_ADDRESS"`
	SecretKey     string `mapstructure:"REDIR_SECRET_KEY"`
}

// NewConfig load the configuration app.
func NewConfig(logger log.Logger, path string) (config Config) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.SetDefault("DBAddr", dbAddr)
	viper.SetDefault("DBPassword", dbPassword)
	viper.SetDefault("DBDatabase", dbDatabase)

	viper.SetDefault("CacheAddr", cacheAddr)
	viper.SetDefault("CachePassword", cachePassword)
	viper.SetDefault("CacheDatabase", cacheDatabase)

	viper.SetDefault("ServerAddress", serverAddress)
	viper.SetDefault("SecretKey", secretKey)

	viper.AutomaticEnv()
	viper.ReadInConfig()

	err := viper.Unmarshal(&config)
	if err != nil {
		logger.Log("method", "NewConfig", "message", "error with viper.Unmarshal", "err", err.Error())
		os.Exit(1)
	}
	return
}
