package middleware

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	"github.com/wvoliveira/corgi/internal/pkg/token"
)

// Authentication check if auth ok and set claims in request header.
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.Logger(c)

		var user model.User
		var err error

		headerAuth := c.GetHeader("Authorization")

		if headerAuth == "" {
			user.ID = "0"
			user.Name = "Anonymous"
			user.Role = "anon"
		}

		headerToken := strings.Split(headerAuth, "Bearer ")
		accessToken := headerToken[len(headerToken)-1]

		if headerAuth != "" && accessToken != "" {
			user, err = token.ValidateToken(accessToken)

			if err != nil {
				log.Error().Caller().Msg(err.Error())
				e.EncodeError(c, e.ErrUnauthorized)
				c.Abort()
			}
		}

		c.Set("user_id", user.ID)
		c.Set("user_role", user.Role)
		c.Next()
	}
}

// Authorization check if auth ok and set claims in request header.
func Authorization(en *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.Logger(c)

		v, ok := c.Get("user_role")
		if !ok {
			log.Error().Caller().Msg("impossible to know who you are")
			e.EncodeError(c, e.ErrUnauthorized)
			c.Abort()
		}

		role := v.(string)

		sub := role                 // role that wants to access a resource.
		obj := c.Request.RequestURI // the resource that is going to be accessed.
		act := c.Request.Method     // the operation that the user performs on the resource.

		log.Debug().Caller().Msg(fmt.Sprintf("Sub: %s, Obj: %s, Act: %s", sub, obj, act))

		ok, _ = en.Enforce(sub, obj, act)
		if !ok {
			// deny the request, show an error
			log.Debug().Caller().Msg("Deny by authorization config")
			e.EncodeError(c, e.ErrUnauthorized)
			c.Abort()
			return
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
				c.Abort()
			}
		}
		c.Next()
	}
}

// UniqueUserForKeywords add session for unique user verification.
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

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()
		log = logger.Logger(c)

		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}

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

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
