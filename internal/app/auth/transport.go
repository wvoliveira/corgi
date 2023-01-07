package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/middleware"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func (s service) NewHTTP(rg *gin.RouterGroup) {
	r := rg.Group("/auth")
	r.Use(middleware.Auth())

	r.GET("/logout", s.HTTPLogout)
}

func (s service) HTTPLogout(c *gin.Context) {
	l := log.Ctx(c)

	user, err := decodeLogout(c)
	if err != nil {
		l.Error().Caller().Msg(err.Error())
		e.EncodeError(c, err)
		return
	}

	err = s.Logout(c, user)
	if err != nil {
		l.Error().Caller().Msg(err.Error())
		e.EncodeError(c, err)
		return
	}

	session := sessions.Default(c)
	session.Delete("user")

	err = session.Save()
	if err != nil {
		l.Error().Caller().Msg(err.Error())
		e.EncodeError(c, err)
		return
	}

	c.JSON(200, response.Response{
		Status:  "successful",
		Message: "Logout success!",
	})
}
