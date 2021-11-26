package password

import (
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/gin-gonic/gin"
)

func (s service) HTTPLogin(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeAuthLoginRequest(c.Request)
	if err != nil {
		return
	}

	// Business logic.
	_, user, err := s.Login(c.Request.Context(), dr.Email, dr.Password)
	if err != nil {
		return
	}

	// Encode object to answer request (response).
	sr := authLoginResponse{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		ExpiresIn:    7200,
		Err:          err,
	}
	_ = encodeResponse(c.Writer, sr)
}

func (s service) HTTPRegister(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeAuthRegisterRequest(c.Request)
	if err != nil {
		e.EncodeError(err, c.Writer)
		return
	}

	// Business logic.
	err = s.Register(c.Request.Context(), dr.Email, dr.Password)
	if err != nil {
		e.EncodeError(err, c.Writer)
		return
	}

	// Encode object to answer request (response).
	sr := authRegisterResponse{Err: err}
	_ = encodeResponse(c.Writer, sr)
}
