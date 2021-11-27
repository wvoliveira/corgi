package token

import (
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/gin-gonic/gin"
)

func (s service) HTTPRefresh(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeRefreshRequest(c.Request)
	if err != nil {
		e.EncodeError(c, err)
		return
	}
	token := entity.Token{AccessToken: dr.AccessToken, RefreshToken: dr.AccessToken, AccessExpires: dr.ExpiresIn}

	// Business logic.
	token, err = s.Refresh(c.Request.Context(), token)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := refreshResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.AccessExpires,
		Err:          err,
	}
	encodeResponse(c, sr)
}
