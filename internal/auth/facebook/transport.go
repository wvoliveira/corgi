package facebook

import (
	"fmt"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/middlewares"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (s service) Routers(e *gin.Engine) {
	r := e.Group("/api/auth/facebook",
		middlewares.Access(s.logger),
		middlewares.Checks(s.logger),
		sessions.SessionsMany([]string{"session_unique", "session_auth"}, s.store),
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
	token, err := s.Callback(c.Request.Context(), callbackURL, dr)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := callbackResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.AccessExpires,
		Err:          err,
	}

	sessionAuth := sessions.DefaultMany(c, "session_auth")
	sessionAuth.Set("access_token", token.AccessToken)
	sessionAuth.Set("refresh_token_id", token.RefreshToken)
	err = sessionAuth.Save()
	if err != nil {
		e.EncodeError(c, err)
		return
	}
	encodeResponse(c, sr)
}
