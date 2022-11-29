package config

import (
	"strings"

	"github.com/mitchellh/go-homedir"
	flag "github.com/spf13/pflag"

	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

// New load the configuration app.
func New(configFile string) {
	viper.SetDefault("app.log_level", "info")
	viper.SetDefault("app.secret_key", "CHANGE_FOR_SOMETHING_MORE_SECURITY")
	viper.SetDefault("app.redirect_url", "http://127.0.0.1:8081")

	viper.SetDefault("server.http_port", 8081)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("CORGI")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName("corgi")
	viper.SetConfigType("yaml")

	file, err := homedir.Expand(configFile)

	if err != nil {
		log.Error().Caller().Msg("error to expand config file: " + err.Error())
	} else {
		viper.SetConfigFile(file)
	}

	err = viper.ReadInConfig()
	if err != nil {
		log.Warn().Caller().Msg(err.Error())
	}

	viper.BindPFlags(flag.CommandLine)
}
