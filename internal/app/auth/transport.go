package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
)

func (s service) NewHTTP(rg *gin.RouterGroup) {
	r := rg.Group("/auth")
	// rr.Use(middleware.Checks)
	// rr.Use(middleware.Auth(s.secret))

	r.GET("/logout", s.HTTPLogout())
}

func (s service) HTTPLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := log.Ctx(c)

		dr, err := decodeLogout(c)
		if err != nil {
			l.Error().Caller().Msg(err.Error())
			e.EncodeError(c, err)
			return
		}

		err = s.Logout(c, dr.Token)
		if err != nil {
			l.Error().Caller().Msg(err.Error())
			e.EncodeError(c, err)
			return
		}

		// TODO: check c.Request.Host if working properly.
		c.SetCookie("token_access", "", -1, "/", c.Request.Host, false, true)
		c.SetCookie("logged", "", -1, "/", c.Request.Host, false, true)

		c.JSON(200, "Logout success!")
	}
}
