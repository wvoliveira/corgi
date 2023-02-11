package config

import (
	"fmt"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

// New load the configuration app.
func New() {
	viper.SetDefault("log_level", "info")
	viper.SetDefault("secret_key", "CHANGE_FOR_SOMETHING_MORE_SECURITY")
	viper.SetDefault("redirect_url", "http://127.0.0.1:8081")
	viper.SetDefault("datasource", "postgres://user:password@localhost:5432/corgi?sslmode=disable")

	viper.SetDefault("server.http_port", 8081)
	viper.SetDefault("server.read_timeout", 10)
	viper.SetDefault("server.write_timeout", 10)

	viper.SetDefault("domain_default", fmt.Sprintf("localhost:%d", viper.GetInt("server.http_port")))
	viper.SetDefault("domain_alternatives", []string{})

	viper.AutomaticEnv()
	viper.SetEnvPrefix("CORGI")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.BindPFlags(flag.CommandLine)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
	}
}
