package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
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
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := entity.User{}

		v := session.Get("user")

		fmt.Println(v)

		if v == nil {
			user.ID = "anonymous"
			user.Name = "Anonymous"

			session.Set("user", user)
			session.Save()
		}

		if v != nil {
			err := json.Unmarshal(v.([]byte), &user)
			if err != nil {
				log.Error().Caller().Msg(err.Error())
			}
		}

		c.Next()
	}
}

// Checks returns a middleware that verify some points before business logic.
func Checks() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Method == "POST" || c.Request.Method == "PATCH" {

			if c.Request.Body == http.NoBody {
				log.Warn().Caller().Msg("Empty body in POST or PATCH request")
				e.EncodeError(c, e.ErrRequestNeedBody)
				return
			}
		}

		c.Next()
	}
}
