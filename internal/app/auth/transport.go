package auth

import (
	"net/http"

	e "github.com/elga-io/corgi/internal/pkg/errors"
	"github.com/elga-io/corgi/internal/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (s service) NewHTTP(e *gin.Engine) {
	r := e.Group("/auth",
		middlewares.Checks(),
		middlewares.Auth(s.secret),
		middlewares.Authorizer(s.enforce))

	r.GET("/logout", s.HTTPLogout)
}

func (s service) HTTPLogout(c *gin.Context) {
	l := log.Ctx(c)

	// Decode request to object.
	dr, err := decodeLogout(c)
	if err != nil {
		l.Error().Caller().Msg(err.Error())
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	err = s.Logout(c.Request.Context(), dr.Token)
	if err != nil {
		l.Error().Caller().Msg(err.Error())
		e.EncodeError(c, err)
		return
	}

	cookieAccess := http.Cookie{Name: "access_token", MaxAge: -1}
	cookieRefresh := http.Cookie{Name: "refresh_token_id", MaxAge: -1}
	cookieLogged := http.Cookie{Name: "logged", MaxAge: -1}

	http.SetCookie(c.Writer, &cookieAccess)
	http.SetCookie(c.Writer, &cookieRefresh)
	http.SetCookie(c.Writer, &cookieLogged)

	// Encode object to answer request (response).
	sr := logoutResponse{Err: err}
	encodeResponse(c, sr)
}
