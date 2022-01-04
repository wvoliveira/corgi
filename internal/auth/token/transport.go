package token

import (
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func (s service) Routers(e *gin.Engine) {
	r := e.Group("/auth/token",
		middlewares.Checks(s.logger),
		middlewares.Auth(s.logger, s.secret),
		middlewares.Authorizer(s.enforce, s.logger))

	r.POST("/refresh", s.HTTPRefresh)
}

func (s service) HTTPRefresh(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeRefreshRequest(c.Request)
	if err != nil {
		e.EncodeError(c, err)
		return
	}
	token := entity.Token{ID: dr.RefreshTokenID}

	// Business logic.
	tokenAccess, tokenRefresh, err := s.Refresh(c.Request.Context(), token)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := refreshResponse{
		AccessToken:  tokenAccess.Token,
		RefreshToken: tokenRefresh.Token,
		ExpiresIn:    tokenAccess.ExpiresIn,
		Err:          err,
	}
	encodeResponse(c, sr)
}
