package middleware

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/elga-io/corgi/internal/app/entity"
	e "github.com/elga-io/corgi/internal/pkg/errors"
	"github.com/elga-io/corgi/internal/pkg/jwt"
	"github.com/elga-io/corgi/internal/pkg/request"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog/log"
)

// Access returns a middleware that records an access log message for every HTTP request being processed.
func Access(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		body, _ := io.Copy(w, r.Body)

		// associate request ID and session ID with the request context
		// so that they can be added to the log messages
		ctx := r.Context()
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

		// End logging response access log.
		log.Info().
			Str("http", "request").
			Str("client_ip", request.IP(r)).
			Int64("duration", time.Since(start).Milliseconds()).
			Int("status", r.Response.StatusCode).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("query", r.URL.RawQuery).
			Str("proto", r.Proto).
			Int64("size", body).
			Str("user-agent", r.UserAgent()).Discard().Send()
	})
}

// Auth check if auth ok and set claims in request header.
func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := log.Ctx(r.Context())

			hashToken, err := r.Cookie("access_token")

			if err == http.ErrNoCookie || hashToken.Value == "" {
				l.Info().Caller().Msg("cookie with the access_token name was not found or blank")
				e.EncodeError(w, e.ErrNoTokenFound)
				return
			}

			// validate ID token here
			claims, err := jwt.ValidateToken(hashToken.Value, secret)

			if err != nil {
				l.Warn().Caller().Msg(fmt.Sprintf("the token is invalid: %s", err.Error()))
				e.EncodeError(w, e.ErrTokenInvalid)
				return
			}

			tokenRefreshID, _ := r.Cookie("refresh_token_id")
			ii := entity.IdentityInfo{
				ID:             claims["identity_id"].(string),
				Provider:       claims["identity_provider"].(string),
				UID:            claims["identity_uid"].(string),
				UserID:         claims["user_id"].(string),
				UserRole:       claims["user_role"].(string),
				RefreshTokenID: tokenRefreshID.Value,
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, entity.IdentityInfo{}, ii)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// Checks returns a middleware that verify some points before business logic.
func Checks(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := log.Ctx(r.Context())

		if r.Method == "POST" || r.Method == "PATCH" {
			if r.Body == http.NoBody {
				l.Warn().Caller().Msg("Empty body in POST or PATCH request")
				e.EncodeError(w, e.ErrRequestNeedBody)
			}
		}
		next.ServeHTTP(w, r)
	})
}

// Authorizer check if user role has access to resource.
func Authorizer(en *casbin.Enforcer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			l := log.Ctx(ctx)

			var (
				role string
				ii   = entity.IdentityInfo{}
			)

			anyy := ctx.Value(ii)
			if anyy != nil {
				ii := anyy.(entity.IdentityInfo)
				role = ii.UserRole
			}

			if role == "" {
				role = "anonymous"
			}

			// casbin rule enforcing
			ok, err := en.Enforce(role, r.URL.Path, r.Method)
			if err != nil {
				l.Error().Caller().Msg(err.Error())
				e.EncodeError(w, err)
				return
			}

			if ok {
				next.ServeHTTP(w, r)
			} else {
				e.EncodeError(w, e.ErrUnauthorized)
				return
			}
		})
	}
}

// SesssionRedirect check if user already clicked in shortener link.
func SesssionRedirect(store *sessions.CookieStore, sessionName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, sessionName)

			next.ServeHTTP(w, r)

			if r.Response.StatusCode == 404 {
				return
			}

			if data := session.Values[r.URL.Path]; data == nil {
				session.Values[r.URL.Path] = true
				err := session.Save(r, w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		})
	}
}
