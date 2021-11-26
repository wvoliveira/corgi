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
			e.EncodeError(e.ErrNoTokenFound, c.Writer)
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
			e.EncodeError(e.ErrTokenExpired, c.Writer)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			accountID := claims["id"].(string)
			accountEmail := claims["email"].(string)
			accountRole := claims["role"].(string)

			c.Set("AccountID", accountID)
			c.Set("AccountEmail", accountEmail)
			c.Set("AccountRole", accountRole)

			c.Next()
		}
		e.EncodeError(e.ErrUnauthorized, c.Writer)
	}
}
