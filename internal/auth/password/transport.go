package password

import (
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (s service) Routers(r *gin.RouterGroup) {
	r.POST("/auth/password/login", s.HTTPLogin)
	r.POST("/auth/password/register", s.HTTPRegister)
	// v1Auth.POST("/google/login", authGoogleService.HTTPLogin)
}

func (s service) HTTPLogin(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeLoginRequest(c.Request)
	if err != nil {
		e.EncodeError(c, err)
		return
	}
	identity := entity.Identity{Provider: "email", UID: dr.Email, Password: dr.Password}

	// Business logic.
	token, err := s.Login(c.Request.Context(), identity)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := loginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.AccessExpires,
		Err:          err,
	}

	sessionAuth := sessions.DefaultMany(c, "session_auth")
	sessionAuth.Set("access_token", token.AccessToken)
	err = sessionAuth.Save()
	if err != nil {
		e.EncodeError(c, err)
		return
	}
	encodeResponse(c, sr)
}

func (s service) HTTPRegister(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeRegisterRequest(c.Request)
	if err != nil {
		e.EncodeError(c, err)
		return
	}
	identity := entity.Identity{Provider: "email", UID: dr.Email, Password: dr.Password}

	// Business logic.
	err = s.Register(c.Request.Context(), identity)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := registerResponse{Err: err}
	encodeResponse(c, sr)
}
