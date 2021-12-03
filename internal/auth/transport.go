package auth

import (
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/middlewares"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (s service) Routers(e *gin.Engine) {
	e.GET("/api/auth/logout",
		s.HTTPLogout,
		sessions.SessionsMany([]string{"session_unique", "session_auth"}, s.store),
		middlewares.Auth(s.logger, s.secret))
}

func (s service) HTTPLogout(c *gin.Context) {
	// Decode request to object.
	dr, err := decodeLogout(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	sessionAuth := sessions.DefaultMany(c, "session_auth")
	dr.Token.ID = sessionAuth.Get("refresh_token_id").(string)

	// Business logic.
	err = s.Logout(c.Request.Context(), dr.Token)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	sessionAuth.Delete("access_token")
	sessionAuth.Delete("refresh_token_id")
	err = sessionAuth.Save()
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := logoutResponse{Err: err}
	encodeResponse(c, sr)
}
