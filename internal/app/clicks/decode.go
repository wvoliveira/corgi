package clicks

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func decodeFind(c *gin.Context) (link string) {
	link = c.Param("link")
	link = strings.Replace(link, "_", "/", -1)
	return
}
