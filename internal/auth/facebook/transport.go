package facebook

import (
	"fmt"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func (s service) Routers(e *gin.Engine) {
	r := e.Group("/auth/facebook",
		middlewares.Checks(s.logger),
		middlewares.Authorizer(s.enforce, s.logger))

	r.GET("/login", s.HTTPLogin)
	r.GET("/callback", s.HTTPCallback)
}

func (s service) HTTPLogin(c *gin.Context) {
	// Decode request to request object.
	_, err := decodeLoginRequest(c.Request)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	schema := "http"
	if c.Request.TLS != nil {
		schema = "https"
	}
	callbackURL := fmt.Sprintf("%s://%s", schema, c.Request.Host+"/api/auth/facebook/callback")
	redirectURL, err := s.Login(c.Request.Context(), callbackURL)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := loginResponse{
		RedirectURL: redirectURL,
		Err:         err,
	}
	if err != nil {
		e.EncodeError(c, err)
	}
	encodeResponse(c, sr)
}

func (s service) HTTPCallback(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeCallbackRequest(c.Request)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	schema := "http"
	if c.Request.TLS != nil {
		schema = "https"
	}
	callbackURL := fmt.Sprintf("%s://%s", schema, c.Request.Host+"/api/auth/facebook/callback")
	tokenAccess, tokenRefresh, err := s.Callback(c.Request.Context(), callbackURL, dr)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := callbackResponse{
		AccessToken:  tokenAccess.Token,
		RefreshToken: tokenRefresh.Token,
		ExpiresIn:    tokenAccess.ExpiresIn,
		Err:          err,
	}

	encodeResponse(c, sr)
}
