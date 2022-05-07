package token

import (
	"github.com/elga-io/corgi/internal/app/entity"
	e "github.com/elga-io/corgi/internal/pkg/errors"
	"github.com/elga-io/corgi/internal/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func (s service) NewHTTP(e *gin.Engine) {
	r := e.Group("/auth/token",
		middlewares.Checks(),
		middlewares.Auth(s.secret),
		middlewares.Authorizer(s.enforce))

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
