package middlewares

import (
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Checks returns a middleware that verify some points before business logic.
func Checks() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Method == "POST" || c.Request.Method == "PATCH" {
			if c.Request.Body == http.NoBody {
				e.EncodeError(c, e.ErrRequestNeedBody)
				return
			}
		}
		c.Next()
	}
}
