package config

import (
	"fmt"
	"os"
	"strings"

	"log"

	"github.com/spf13/viper"
)

const (
	appSecretKey     = "changeme"
	appLogLevel      = "info"
	appAdminPassword = "12345"
	appUserPassword  = "12345"
	appRedirectURL   = "http://localhost:4200"

	serverHTTPPort = "8081"
	serverGRPCPort = "8082"

	dbHost     = "127.0.0.1"
	dbPort     = 26257
	dbUser     = "root"
	dbPassword = ""
	dbDatabase = "corgi"

	cacheHost     = "127.0.0.1"
	cachePort     = 6379
	cacheUser     = ""
	cachePassword = ""
	cacheDatabase = 0

	searchHost     = "127.0.0.1"
	searchPort     = 9200
	searchUser     = ""
	searchPassword = ""

	brokerHost     = "127.0.0.1"
	brokerPort     = 4222
	brokerUser     = ""
	brokerPassword = ""
)

// Config a struct for app configuration.
type Config struct {
	App struct {
		SecretKey     string `mapstructure:"secret_key"`
		LogLevel      string `mapstructure:"log_level"`
		AdminPassword string `mapstructure:"admin_password"`
		UserPassword  string `mapstructure:"user_password"`
		RedirectURL   string `mapstructure:"redirect_url"`
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
		Host     string `mapstructure:"host" mapstructure:"DATABASE_HOST"`
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

	Broker struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
	}
}

// NewConfig load the configuration app.
func NewConfig(path string) (config Config) {
	conf := viper.New()

	conf.AutomaticEnv()
	conf.SetEnvPrefix("CORGI")
	conf.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	conf.SetConfigName("app")
	conf.SetConfigType("yaml")
	conf.AddConfigPath(path)

	err := conf.ReadInConfig()
	if err != nil {
		fmt.Printf("Error to read config file: %s\n", err.Error())
	}

	conf.SetConfigFile(".env.yaml")
	conf.AddConfigPath(".")
	err = conf.MergeInConfig()
	if err != nil {
		fmt.Printf("Fatal error config file: %s\n", err.Error())
	}

	// App config.
	conf.SetDefault("app.secret_key", appSecretKey)
	conf.SetDefault("app.log_level", appLogLevel)
	conf.SetDefault("app.admin_password", appAdminPassword)
	conf.SetDefault("app.user_password", appUserPassword)
	conf.SetDefault("app.redirect_url", appRedirectURL)

	// Server config.
	conf.SetDefault("server.http_port", serverHTTPPort)
	conf.SetDefault("server.grpc_port", serverGRPCPort)

	// Database config.
	conf.SetDefault("database.host", dbHost)
	conf.SetDefault("database.port", dbPort)
	conf.SetDefault("database.user", dbUser)
	conf.SetDefault("database.password", dbPassword)
	conf.SetDefault("database.database", dbDatabase)

	// Cache config.
	conf.SetDefault("cache.host", cacheHost)
	conf.SetDefault("cache.port", cachePort)
	conf.SetDefault("cache.user", cacheUser)
	conf.SetDefault("cache.password", cachePassword)
	conf.SetDefault("cache.database", cacheDatabase)

	// Search engine config.
	conf.SetDefault("search.host", searchHost)
	conf.SetDefault("search.port", searchPort)
	conf.SetDefault("search.user", searchUser)
	conf.SetDefault("search.password", searchPassword)

	// Broker config.
	conf.SetDefault("broker.host", brokerHost)
	conf.SetDefault("broker.port", brokerPort)
	conf.SetDefault("broker.user", brokerUser)
	conf.SetDefault("broker.password", brokerPassword)

	// workaround because viper does not treat env vars the same as other config
	for _, key := range conf.AllKeys() {
		val := conf.Get(key)
		conf.Set(key, val)
	}

	if err := conf.Unmarshal(&config); err != nil {
		log.Println("error to load config", "method", "NewConfig", "err", err.Error())
		os.Exit(1)
	}
	return
}
