package password

import (
	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func encodeRegister(c *gin.Context) {
	c.JSON(200, response.Response{Status: "successful"})
}
