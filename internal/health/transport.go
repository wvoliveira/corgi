package health

import (
	"github.com/elga-io/corgi/pkg/middlewares"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (s service) Routers(e *gin.Engine) {
	e.GET("/health",
		s.httpHealth,
		sessions.SessionsMany([]string{"session_unique", "session_auth"}, s.store),
		middlewares.Authorizer(s.enforce, s.logger))
}

func (s service) httpHealth(c *gin.Context) {
	c.JSON(200, "OK "+s.version)
}
