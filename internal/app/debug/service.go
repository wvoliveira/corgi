package debug

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Config() map[string]interface{}
	Env() map[string]string

	NewHTTP(*gin.RouterGroup)
	HTTPConfig(*gin.Context)
	HTTPEnv(*gin.Context)
}

type service struct{}

// NewService creates a new debug "service".
func NewService() Service {
	return service{}
}

// Config get info from config struct.
func (s service) Config() (info map[string]interface{}) {
	keys := viper.AllKeys()
	info = make(map[string]interface{})

	for _, key := range keys {
		value := viper.Get(key)
		info[key] = value
	}
	return
}

// Env get info from environment variables.
func (s service) Env() (info map[string]string) {
	vars := os.Environ()
	info = make(map[string]string)

	for _, v := range vars {
		values := strings.Split(v, "=")
		info[values[0]] = values[1]
	}
	return
}
