package auth

import (
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/middlewares"
	"github.com/gin-gonic/gin"
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

	// Encode object to answer request (response).
	sr := logoutResponse{Err: err}
	encodeResponse(c, sr)
}
