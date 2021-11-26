package password

import (
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/gin-gonic/gin"
)

func (s service) HTTPLogin(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeAuthLoginRequest(c.Request)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	token, err := s.Login(c.Request.Context(), dr.Email, dr.Password)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := authLoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.AtExpires,
		Err:          err,
	}
	encodeResponse(c, sr)
}

func (s service) HTTPRegister(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeAuthRegisterRequest(c.Request)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	err = s.Register(c.Request.Context(), dr.Email, dr.Password)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := authRegisterResponse{Err: err}
	encodeResponse(c, sr)
}
