package auth

import (
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/middlewares"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (s service) Routers(e *gin.Engine) {
	r := e.Group("/api/auth",
		middlewares.Access(s.logger),
		middlewares.Checks(s.logger),
		sessions.SessionsMany([]string{"session_unique", "session_auth"}, s.store),
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

	sessionAuth := sessions.DefaultMany(c, "session_auth")
	tokenIDInterface := sessionAuth.Get("refresh_token_id")
	if tokenIDInterface != nil {
		dr.Token.ID = tokenIDInterface.(string)
	} else {
		logger.Error("impossible to get refresh_token_id from session")
	}

	// Business logic.
	err = s.Logout(c.Request.Context(), dr.Token)
	if err != nil {
		logger.Error("error in service logout: ", err.Error())
		e.EncodeError(c, err)
		return
	}

	sessionAuth.Delete("access_token")
	sessionAuth.Delete("refresh_token_id")
	err = sessionAuth.Save()
	if err != nil {
		logger.Error("error in session auth save: ", err.Error())
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := logoutResponse{Err: err}
	encodeResponse(c, sr)
}
