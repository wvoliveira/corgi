package config

import (
	"os"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

const (
	secretKey     = "changeme"
	logLevel      = "info"
	adminPassword = "12345"
	userPassword  = "12345"

	redirectURL = "http://127.0.0.1:8081"

	serverHTTPPort = "8081"
)

// Config a struct for app configuration.
type Config struct {
	SecretKey     string `mapstructure:"secret_key"`
	LogLevel      string `mapstructure:"log_level"`
	AdminPassword string `mapstructure:"admin_password"`
	UserPassword  string `mapstructure:"user_password"`

	RedirectURL string `mapstructure:"redirect_url"`

	Server struct {
		HTTPPort string `mapstructure:"http_port"`
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
}

// NewConfig load the configuration app.
func NewConfig(configFile string) (config Config) {
	conf := viper.New()

	conf.AutomaticEnv()
	conf.SetEnvPrefix("CORGI")
	conf.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	conf.SetConfigFile(configFile)
	conf.SetConfigType("yaml")
	err := conf.ReadInConfig()
	if err != nil {
		log.Warn().Caller().Msg(err.Error())
	}

	conf.SetDefault("secret_key", secretKey)
	conf.SetDefault("log_level", logLevel)
	conf.SetDefault("admin_password", adminPassword)
	conf.SetDefault("user_password", userPassword)

	conf.SetDefault("redirect_url", redirectURL)

	conf.SetDefault("server.http_port", serverHTTPPort)

	// workaround because viper does not treat env vars the same as other config
	for _, key := range conf.AllKeys() {
		val := conf.Get(key)
		conf.Set(key, val)
	}

	if err := conf.Unmarshal(&config); err != nil {
		log.Fatal().Caller().Msg(err.Error())
		os.Exit(1)
	}
	return
}
