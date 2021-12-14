package health

import (
	"github.com/elga-io/corgi/pkg/middlewares"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (s service) Routers(e *gin.Engine) {
	r := e.Group("/health",
		sessions.SessionsMany([]string{"unique", "auth"}, s.store),
		middlewares.Authorizer(s.enforce, s.logger))
	r.GET("/ping", s.httpHealth)
	//r.GET("/live", s.httpLive)
	//r.GET("/ready", s.httpReady)
}

func (s service) httpHealth(c *gin.Context) {
	c.JSON(200, "pong "+s.version)
}
