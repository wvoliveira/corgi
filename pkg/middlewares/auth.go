package middlewares

import (
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
)

// Auth check if auth ok and set claims in request header.
func Auth(logger log.Logger, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		logg := logger.With(c.Request.Context())

		sessionAuth := sessions.DefaultMany(c, "session_auth")
		if sessionAuth == nil {
			logg.Info("session_auth not found")
			_ = c.AbortWithError(http.StatusUnauthorized, e.ErrNoTokenFound)
			return
		}

		tokenInterface := sessionAuth.Get("access_token")
		if tokenInterface == nil {
			logg.Info("access token not found in session cookies")
			_ = c.AbortWithError(http.StatusUnauthorized, e.ErrNoTokenFound)
			return
		}

		accessToken := tokenInterface.(string)
		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logg.Warnf("fail to parse access token")
				_ = c.AbortWithError(http.StatusUnauthorized, e.ErrTokenInvalid)
				return token, e.ErrParseToken
			}
			return []byte(secret), nil
		})

		if err != nil {
			logg.Infof("error to parse access token: %s", err.Error())
			_ = c.AbortWithError(http.StatusUnauthorized, e.ErrTokenExpired)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("identity_id", claims["identity_id"].(string))
			c.Set("identity_provider", claims["identity_provider"].(string))
			c.Set("identity_uid", claims["identity_uid"].(string))
			c.Set("user_id", claims["user_id"].(string))
			c.Set("user_role", claims["user_role"].(string))
			c.Next()
		} else {
			logg.Info("invalid token! so sorry")
			_ = c.AbortWithError(http.StatusUnauthorized, e.ErrTokenInvalid)
		}
		return
	}
}
