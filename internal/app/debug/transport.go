package debug

import (
	"net/http"

	"github.com/gin-gonic/gin"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func (s service) NewHTTP(rg *gin.RouterGroup) {
	r := rg.Group("/debug")

	r.GET("/config", s.HTTPConfig)
	r.GET("/env", s.HTTPEnv)
	//r.GET("/live", s.httpLive)
	//r.GET("/ready", s.httpReady)
}

func (s service) HTTPConfig(c *gin.Context) {
	data, err := s.Info(r.Context())
	if err != nil {
		e.EncodeError(w, err)
		return
	}
	response.Default(w, data, "", http.StatusOK)
}

func (s service) HTTPEnv(c *gin.Context) {
	data, err := s.Info(r.Context())
	if err != nil {
		e.EncodeError(w, err)
		return
	}
	response.Default(w, data, "", http.StatusOK)
}
