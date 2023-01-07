package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func (s service) NewHTTP(rg *gin.RouterGroup) {
	r := rg.Group("/health")
	r.GET("", s.HTTPHealth)
	r.GET("/live", s.HTTPLive)
	r.GET("/ready", s.HTTPReady)
}

func (s service) HTTPHealth(c *gin.Context) {

	healths, err := s.Health(c)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	httpStatusCode := http.StatusOK

	for _, item := range healths {
		if item.Required && item.Status != "OK" {
			httpStatusCode = http.StatusServiceUnavailable
		}
	}

	response.Default(c, healths, "", httpStatusCode)
}

func (s service) HTTPLive(c *gin.Context) {
	response.Default(c, "Live", "", http.StatusOK)
}

func (s service) HTTPReady(c *gin.Context) {
	response.Default(c, "Ready", "", http.StatusOK)
}
