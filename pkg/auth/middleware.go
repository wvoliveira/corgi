package auth

import (
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// MiddlewareAuth check if auth ok and set claims in request header.
func MiddlewareAuth(logger log.Logger, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header["Token"] == nil {
			e.EncodeError(c, e.ErrNoTokenFound)
			return
		}

		var SigningKey = []byte(secret)

		token, err := jwt.Parse(c.Request.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, e.ErrParseToken
			}
			return SigningKey, nil
		})

		if err != nil {
			e.EncodeError(c, e.ErrTokenExpired)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("identity_id", claims["identity_id"].(string))
			c.Set("identity_provider", claims["identity_provider"].(string))
			c.Set("identity_uid", claims["identity_uid"].(string))
			c.Set("user_id", claims["user_id"].(string))
			c.Set("user_role", claims["user_role"].(string))
			c.Next()
		}
		e.EncodeError(c, e.ErrUnauthorized)
	}
}
