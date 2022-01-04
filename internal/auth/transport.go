package auth

import (
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s service) Routers(e *gin.Engine) {
	r := e.Group("/auth",
		middlewares.Checks(s.logger),
		middlewares.Auth(s.logger, s.secret),
		middlewares.Authorizer(s.enforce, s.logger))

	r.GET("/logout", s.HTTPLogout)
}

func (s service) HTTPLogout(c *gin.Context) {
	logger := s.logger.With(c)

	// Decode request to object.
	dr, err := decodeLogout(c)
	if err != nil {
		logger.Error("error in decode logout: ", err.Error())
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	err = s.Logout(c.Request.Context(), dr.Token)
	if err != nil {
		logger.Error("error in service logout: ", err.Error())
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
