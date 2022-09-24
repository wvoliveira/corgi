package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog/log"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/jwt"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/request"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Access returns a middleware that records an access log message for every HTTP request being processed.
func Access(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Disable logs for some paths.
		if strings.HasPrefix(r.URL.Path, "/_next") {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		// Copy body payload to get length.
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
		}

		// Insert body again to use in another handlers.
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		lrw := NewLoggingResponseWriter(w)

		// associate request ID and session ID with the request context
		// so that they can be added to the log messages
		ctx := r.Context()
		if v := ctx.Value(entity.CorrelationID{}); v == nil {
			ctx = context.WithValue(ctx, entity.CorrelationID{}, entity.CorrelationID{ID: uuid.New().String()})
		}

		r = r.WithContext(ctx)
		next.ServeHTTP(lrw, r)

		l := logger.Logger(ctx)
		l.Info().
			Caller().
			Str("client_ip", request.IP(r)).
			Float64("duration", time.Since(start).Seconds()).
			Int("status", lrw.statusCode).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("query", r.URL.RawQuery).
			Str("proto", r.Proto).
			Int("size", len(body)).
			Str("user-agent", r.UserAgent()).Msg("request")
	})
}

// Auth check if auth ok and set claims in request header.
func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := logger.Logger(r.Context())

			// In this app, we can create link without authentication.
			// So, in some routes we can forward without user_id.
			// But only for period of time, like expiration in database/cache.
			anonymousAccess := false

			hashToken, err := r.Cookie("access_token")
			if err == http.ErrNoCookie || hashToken.Value == "" {
				anonymousAccess = true
			}

			identity := entity.IdentityInfo{}

			if !anonymousAccess {
				// validate ID token here
				claims, err := jwt.ValidateToken(hashToken.Value, secret)

				if err != nil {
					l.Warn().Caller().Msg(fmt.Sprintf("the token is invalid: %s", err.Error()))
					e.EncodeError(w, e.ErrTokenInvalid)
					return
				}

				tokenRefreshID, err := r.Cookie("refresh_token_id")
				if err != nil {
					l.Warn().Caller().Msg(fmt.Sprintf("error to get refresh_token_id from cookie: %s", err.Error()))
					e.EncodeError(w, e.ErrTokenInvalid)
					return
				}

				identity = entity.IdentityInfo{
					ID:             claims["identity_id"].(string),
					Provider:       claims["identity_provider"].(string),
					UID:            claims["identity_uid"].(string),
					UserID:         claims["user_id"].(string),
					UserRole:       claims["user_role"].(string),
					RefreshTokenID: tokenRefreshID.Value,
				}
			}

			if anonymousAccess {
				identity = entity.IdentityInfo{
					UserID: "anonymous",
				}
			}

			// Authorizer. Casbin was removed but I'm rethinking if its was a good idea.
			// Because now I need to check if user is anonymous access and block for some paths manually.
			// Oh my lord.
			if identity.UserID == "anonymous" {
				path := r.URL.Path
				unauthorizedPaths := []string{"/api/v1/groups", "/api/v1/auth/token", "/api/v1/auth/logout"}

				for _, unauthorizedPath := range unauthorizedPaths {
					if strings.HasPrefix(path, unauthorizedPath) {
						err = e.ErrUnauthorized
						e.EncodeError(w, err)
						return
					}
				}
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, entity.IdentityInfo{}, identity)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// Checks returns a middleware that verify some points before business logic.
// Like POST and PATCH without request body.
func Checks(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PATCH" {
			if r.Body == http.NoBody {
				log.Warn().Caller().Msg("Empty body in POST or PATCH request")
				e.EncodeError(w, e.ErrRequestNeedBody)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// SessionRedirect check if user already clicked in shortener link.
func SessionRedirect(store *sessions.CookieStore, sessionName string) func(http.Handler) http.Handler {
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
