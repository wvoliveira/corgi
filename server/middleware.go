package server

import (
	"net/http"

	"github.com/golang-jwt/jwt"
)

/*
  Access control for Web UI.
*/

// AccessControl set common headers for web UI.
func AccessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

// IsAuthorized check if auth ok and set claims in request header.
func IsAuthorized(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			encodeError(ErrNoTokenFound, w)
			return
		}

		var mySigningKey = []byte(secretKey)

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

			h.ServeHTTP(w, r)
			return
		}
		encodeError(ErrUnauthorized, w)
	}
}
