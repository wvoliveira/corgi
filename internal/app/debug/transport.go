package debug

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func (s service) NewHTTP(rg *gin.RouterGroup) {
	r := rg.Group("/debug")

	r.GET("/config", s.HTTPConfig)
	r.GET("/env", s.HTTPEnv)
}

func (s service) HTTPConfig(c *gin.Context) {
	data := s.Config()
	response.Default(c, data, "", http.StatusOK)
}

func (s service) HTTPEnv(c *gin.Context) {
	data := s.Env()
	response.Default(c, data, "", http.StatusOK)
}
