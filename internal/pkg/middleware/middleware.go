package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/casbin/casbin/v2"
	e "github.com/elga-io/corgi/internal/pkg/errors"
	"github.com/elga-io/corgi/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type authHeader struct {
	IDToken string `header:"Authorization"`
}

// Access returns a middleware that records an access log message for every HTTP request being processed.
func Access() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// associate request ID and session ID with the request context
		// so that they can be added to the log messages
		ctx := c.Request.Context()
		c.Request = c.Request.WithContext(ctx)
		c.Next()

		// End logging response access log.
		log.Info().
			Str("http", "request").
			Str("client_ip", c.ClientIP()).
			Int64("duration", time.Since(start).Milliseconds()).
			Int("status", c.Writer.Status()).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("query", c.Request.URL.RawQuery).
			Str("proto", c.Request.Proto).
			Int("status", c.Writer.Status()).
			Int("size", c.Writer.Size()).
			Str("user-agent", c.Request.UserAgent()).Discard().Send()
	}
}

// Auth check if auth ok and set claims in request header.
func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		l := log.Ctx(c.Request.Context())

		hashToken, err := c.Cookie("access_token")

		if err == http.ErrNoCookie || hashToken == "" {
			l.Info().Caller().Msg("cookie with the access_token name was not found or blank")
			_ = c.AbortWithError(http.StatusUnauthorized, e.ErrNoTokenFound)
			return
		}

		// validate ID token here
		claims, err := jwt.ValidateToken(hashToken, secret)

		if err != nil {
			l.Warn().Caller().Msg(fmt.Sprintf("the token is invalid: %s", err.Error()))
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
func Checks() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := log.Ctx(c.Request.Context())

		if c.Request.Method == "POST" || c.Request.Method == "PATCH" {
			if c.Request.Body == http.NoBody {
				l.Warn().Caller().Msg("Empty body in POST or PATCH request")
				_ = c.AbortWithError(http.StatusBadRequest, e.ErrRequestNeedBody)
			}
		}
		c.Next()
	}
}

// Authorizer check if user role has access to resource.
func Authorizer(en *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		l := log.Ctx(c.Request.Context())

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
			l.Error().Caller().Msg(err.Error())
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
