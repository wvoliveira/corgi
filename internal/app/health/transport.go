package health

import (
	"github.com/elga-io/corgi/internal/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func (s service) NewHTTP(e *gin.Engine) {
	r := e.Group("/health",
		middlewares.Authorizer(s.enforce))

	r.GET("/ping", s.HTTPHealth)
	//r.GET("/live", s.httpLive)
	//r.GET("/ready", s.httpReady)
}

func (s service) HTTPHealth(c *gin.Context) {
	c.JSON(200, "pong "+s.version)
}
