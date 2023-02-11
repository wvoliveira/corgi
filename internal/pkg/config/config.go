package config

import (
	"github.com/spf13/viper"
)

// New load the configuration app.
func New() {
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("SECRET_KEY", "CHANGE_FOR_SOMETHING_MORE_SECURITY")
	viper.SetDefault("REDIRECT_URL", "http://127.0.0.1:8081")

	viper.SetDefault("DB_URL", "postgres://user:password@localhost:5432/corgi?sslmode=disable")
	viper.SetDefault("CACHE_URL", "redis://password@localhost:6379/0")

	viper.SetDefault("SERVER_HTTP_PORT", 8081)
	viper.SetDefault("SERVER_READ_TIMEOUT", 10)
	viper.SetDefault("SERVER_WRITE_TIMEOUT", 10)

	viper.SetDefault("DOMAIN_DEFAULT", "localhost:8081")
	viper.SetDefault("DOMAIN_ALTERNATIVES", []string{})

	viper.SetEnvPrefix("CORGI")
	viper.AutomaticEnv()
}
