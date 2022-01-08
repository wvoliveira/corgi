package config

import (
	"fmt"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/spf13/viper"
	"os"
)

const (
	appSecretKey     = "changeme"
	appLogLevel      = "info"
	appAdminPassword = "12345"
	appUserPassword  = "12345"

	serverHTTPPort = "8081"
	serverGRPCPort = "8082"

	dbHost     = "localhost"
	dbPort     = 26257
	dbUser     = "root"
	dbPassword = ""
	dbDatabase = "corgi"

	cacheHost     = "localhost"
	cachePort     = 6379
	cacheUser     = ""
	cachePassword = ""
	cacheDatabase = 0

	searchHost     = "localhost"
	searchPort     = 9200
	searchUser     = ""
	searchPassword = ""
)

// Config a struct for app configuration.
type Config struct {
	App struct {
		SecretKey     string `mapstructure:"secret_key"`
		LogLevel      string `mapstructure:"log_level"`
		AdminPassword string `mapstructure:"admin_password"`
		UserPassword  string `mapstructure:"user_password"`
	}

	Auth struct {
		Google struct {
			ClientID     string `mapstructure:"client_id"`
			ClientSecret string `mapstructure:"client_secret"`
		}
		Facebook struct {
			ClientID     string `mapstructure:"client_id"`
			ClientSecret string `mapstructure:"client_secret"`
		}
	}

	Server struct {
		HTTPPort string `mapstructure:"http_port"`
		GRCPPort string `mapstructure:"grpc_port"`
	}

	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Base     string `mapstructure:"database"`
	}

	Cache struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Database int    `mapstructure:"database"`
	}

	Search struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
	}
}

// NewConfig load the configuration app.
func NewConfig(logger log.Logger, path string) (config Config) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CORGI_")

	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w\n", err))
	}

	viper.SetConfigFile(".env.yaml")
	viper.AddConfigPath(".")
	err = viper.MergeInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w\n", err))
	}

	// App config.
	viper.SetDefault("app.secret_key", appSecretKey)
	viper.SetDefault("app.log_level", appLogLevel)
	viper.SetDefault("app.admin_password", appAdminPassword)
	viper.SetDefault("app.user_password", appUserPassword)

	// Server config.
	viper.SetDefault("server.http_port", serverHTTPPort)
	viper.SetDefault("server.grpc_port", serverGRPCPort)

	// Database config.
	viper.SetDefault("database.host", dbHost)
	viper.SetDefault("database.port", dbPort)
	viper.SetDefault("database.user", dbUser)
	viper.SetDefault("database.password", dbPassword)
	viper.SetDefault("database.database", dbDatabase)

	// Cache config.
	viper.SetDefault("cache.host", cacheHost)
	viper.SetDefault("cache.port", cachePort)
	viper.SetDefault("cache.user", cacheUser)
	viper.SetDefault("cache.password", cachePassword)
	viper.SetDefault("cache.database", cacheDatabase)

	// Search engine config.
	viper.SetDefault("search.host", searchHost)
	viper.SetDefault("search.port", searchPort)
	viper.SetDefault("search.user", searchUser)
	viper.SetDefault("search.password", searchPassword)

	if err := viper.Unmarshal(&config); err != nil {
		logger.Errorf("error to load config", "method", "NewConfig", "err", err.Error())
		os.Exit(1)
	}
	return
}
