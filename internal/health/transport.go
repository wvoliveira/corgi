package health

import (
	"github.com/gin-gonic/gin"
)

func (s service) Routers(r *gin.RouterGroup) {
	r.GET("/health", s.httpHealth)
}

func (s service) httpHealth(c *gin.Context) {
	c.JSON(200, "OK "+s.version)
}
