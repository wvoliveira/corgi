// Package accesslog provides a middleware that records every RESTful API call in a log message.
package middlewares

import (
	"github.com/elga-io/corgi/pkg/log"
	"github.com/gin-gonic/gin"
	"time"
)

// Access returns a middleware that records an access log message for every HTTP request being processed.
func Access(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// associate request ID and session ID with the request context
		// so that they can be added to the log messages
		ctx := c.Request.Context()
		ctx = log.WithRequest(ctx, c.Request)
		c.Request = c.Request.WithContext(ctx)

		// Start logging request access log.
		logger.With(ctx, "http", "request", "start", start).
			Infof("%s %s %s", c.Request.Method, c.Request.URL.Path, c.Request.Proto)

		c.Next()

		// End logging response access log.
		logger.With(ctx, "http", "response", "duration", time.Now().Sub(start).Milliseconds(), "status", c.Writer.Status()).
			Infof("%s %s %s %d %d", c.Request.Method, c.Request.URL.Path, c.Request.Proto, c.Writer.Status(), c.Writer.Size())
	}
}
