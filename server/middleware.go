package server

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type Middlewares struct {
	logger zap.SugaredLogger
	config Config
}

// NewService create a new service with database and cache.
func NewMiddlewares(logger zap.SugaredLogger, config Config) Middlewares {
	return Middlewares{
		logger: logger,
		config: config,
	}
}

// AccessControl set common headers for web UI.
func (m Middlewares) AccessControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

// IsAuthorized check if auth ok and set claims in request header.
func (m Middlewares) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			encodeError(ErrNoTokenFound, w)
			return
		}

		var mySigningKey = []byte(m.config.SecretKey)

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrParseToken
			}
			return mySigningKey, nil
		})

		if err != nil {
			encodeError(ErrTokenExpired, w)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			accountID := claims["id"].(string)
			accountEmail := claims["email"].(string)
			accountRole := claims["role"].(string)

			r.Header.Set("AccountID", accountID)
			r.Header.Set("AccountEmail", accountEmail)
			r.Header.Set("AccountRole", accountRole)

			next.ServeHTTP(w, r)
			return
		}
		encodeError(ErrUnauthorized, w)
	})
}
