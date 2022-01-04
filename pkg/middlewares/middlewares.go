package middlewares

import (
	"github.com/casbin/casbin/v2"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/jwt"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type authHeader struct {
	IDToken string `header:"Authorization"`
}

// Access returns a middleware that records an access log message for every HTTP request being processed.
func Access(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// associate request ID and session ID with the request context
		// so that they can be added to the log messages
		ctx := c.Request.Context()
		ctx = log.WithRequest(ctx, c.Request)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		// End logging response access log.
		logger.With(ctx,
			"http", "request",
			"client_ip", c.ClientIP(),
			"duration", time.Since(start).Milliseconds(),
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
			"proto", c.Request.Proto,
			"status", c.Writer.Status(),
			"size", c.Writer.Size(),
		).Info()
	}
}

// Auth check if auth ok and set claims in request header.
func Auth(logger log.Logger, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		logg := logger.With(c.Request.Context())

		h := authHeader{}

		// bind Authorization Header to h and check for validation errors
		if err := c.ShouldBindHeader(&h); err != nil {
			logg.Info("access token not found Headers")
			_ = c.AbortWithError(http.StatusUnauthorized, e.ErrNoTokenFound)
			return
		}

		hashTokenHeader := strings.Split(h.IDToken, "Bearer ")

		if len(hashTokenHeader) < 2 {
			logg.Info("the Authorization header was found but weird format")
			_ = c.AbortWithError(http.StatusUnauthorized, e.ErrAuthHeaderFormat)
			return
		}

		// validate ID token here
		claims, err := jwt.ValidateToken(hashTokenHeader[1], secret)

		if err != nil {
			logg.Warnf("the token is invalid: %s", err.Error())
			_ = c.AbortWithError(http.StatusUnauthorized, e.ErrTokenInvalid)
			return
		}

		c.Set("identity_id", claims["identity_id"].(string))
		c.Set("identity_provider", claims["identity_provider"].(string))
		c.Set("identity_uid", claims["identity_uid"].(string))
		c.Set("user_id", claims["user_id"].(string))
		c.Set("user_role", claims["user_role"].(string))

		tokenRefreshID, _ := c.Cookie("refresh_token_id")
		c.Set("refresh_token_id", tokenRefreshID)
		c.Next()
	}
}

// Checks returns a middleware that verify some points before business logic.
func Checks(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logg := logger.With(c.Request.Context())

		if c.Request.Method == "POST" || c.Request.Method == "PATCH" {
			if c.Request.Body == http.NoBody {
				logg.Warnf("Empty body in POST or PATCH request")
				_ = c.AbortWithError(http.StatusBadRequest, e.ErrRequestNeedBody)
			}
		}
		c.Next()
	}
}

// Authorizer check if user role has access to resource.
func Authorizer(en *casbin.Enforcer, logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logg := logger.With(c.Request.Context())

		var role string
		roleInterface, ok := c.Get("user_role")
		if ok {
			role = roleInterface.(string)
		}

		if role == "" {
			role = "anonymous"
		}

		// casbin rule enforcing
		res, err := en.Enforce(role, c.Request.URL.Path, c.Request.Method)
		if err != nil {
			logg.Error("error to enforce casbin authorization: ", err.Error())
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if res {
			c.Next()
		} else {
			_ = c.AbortWithError(http.StatusForbidden, e.ErrUnauthorized)
			return
		}
	}
}
