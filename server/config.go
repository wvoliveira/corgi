package server

import (
	"os"

	"github.com/go-kit/log"
	"github.com/spf13/viper"
)

const (
	dbType   = "persistent"
	dbDriver = "sqlite"
	dbSource = "redir.db"

	serverAddress = "0.0.0.0:8080"
	secretKey     = "changeme"
)

// Config a struct for app configuration.
type Config struct {
	DBType   string `mapstructure:"REDIR_DB_TYPE"`
	DBDriver string `mapstructure:"REDIR_DB_DRIVER"`
	DBSource string `mapstructure:"REDIR_DB_SOURCE"`

	ServerAddress string `mapstructure:"REDIR_SERVER_ADDRESS"`
	SecretKey     string `mapstructure:"REDIR_SECRET_KEY"`
}

// NewConfig load the configuration app.
func NewConfig(logger log.Logger, path string) (config Config) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		logger.Log("method", "NewConfig", "message", "error with viper.ReadInConfig", "err", err.Error())
		os.Exit(1)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		logger.Log("method", "NewConfig", "message", "error with viper.Unmarshal", "err", err.Error())
		os.Exit(1)
	}
	return
}
