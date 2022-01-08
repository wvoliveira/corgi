package password

import (
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s service) Routers(e *gin.Engine) {
	r := e.Group("/auth/password",
		middlewares.Checks(s.logger))

	r.POST("/login", s.HTTPLogin)
	r.POST("/register", s.HTTPRegister)
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
	tokenAccess, tokenRefresh, err := s.Login(c.Request.Context(), identity)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	//sr := loginResponse{
	//	AccessToken:  tokenAccess.Token,
	//	RefreshToken: tokenRefresh.Token,
	//	ExpiresIn:    tokenAccess.ExpiresIn,
	//	Err:          err,
	//}

	cookieAccess := http.Cookie{
		Name:    "access_token",
		Value:   tokenAccess.Token,
		Path:    "/",
		Expires: tokenAccess.ExpiresIn,
		// RawExpires
		Secure:   false,
		HttpOnly: false,
	}

	cookieRefresh := http.Cookie{
		Name:    "refresh_token_id",
		Value:   tokenRefresh.ID,
		Path:    "/",
		Expires: tokenRefresh.ExpiresIn,
		// RawExpires
		Secure:   false,
		HttpOnly: false,
	}

	http.SetCookie(c.Writer, &cookieAccess)
	http.SetCookie(c.Writer, &cookieRefresh)

	c.Writer.WriteHeader(200)
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
