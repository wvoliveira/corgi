package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	"github.com/wvoliveira/corgi/internal/pkg/token"
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

// Auth check if auth ok and set claims in request header.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.Logger(c)

		user := model.User{}

		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			user.ID = "0"
			user.Name = "Anonymous"
		}

		if tokenString != "" {
			user, err := token.ValidateToken(tokenString)

			if err != nil {
				log.Error().Caller().Msg(err.Error())
				e.EncodeError(c, e.ErrInternalServerError)
				c.Abort()
			}

			log.Info().Caller().Msg(user.ID)
		}

		c.Set("user_id", user.ID)

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
				c.Abort()
			}
		}

		c.Next()
	}
}

// UniqueUser add session for unique user verification.
func UniqueUserForKeywords() gin.HandlerFunc {
	return func(c *gin.Context) {

		var keywords []string
		var keyword = c.Param("keyword")

		session := sessions.Default(c)
		v := session.Get("keywords")

		if v == nil {
			session.Set("keywords", []string{keyword})
			err := session.Save()

			if err != nil {
				log.Error().Caller().Msg(err.Error())
				e.EncodeError(c, e.ErrRequestNeedBody)
				return
			}

			c.Next()
			return
		}

		keywords = v.([]string)

		for _, k := range keywords {

			if k == keyword {
				c.Next()
				return
			}

		}

		keywords = append(keywords, keyword)

		session.Set("keywords", keywords)
		err := session.Save()

		if err != nil {
			log.Error().Caller().Msg(err.Error())
			e.EncodeError(c, e.ErrRequestNeedBody)
			return
		}

		c.Next()
	}
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.Logger(c)

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Update context with more info.
		log = logger.Logger(c)

		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}

		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

		param.BodySize = c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		param.Path = path

		log.Info().Caller().
			Time("timestamp", param.TimeStamp).
			Str("client_ip", param.ClientIP).
			Str("method", param.Method).
			Str("path", param.Path).
			Str("proto", param.Request.Proto).
			Int("status_code", param.StatusCode).
			Dur("latency", param.Latency).
			Str("user_agent", param.Request.UserAgent()).
			Msg(param.ErrorMessage)
	}
}
