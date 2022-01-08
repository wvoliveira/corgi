package google

import (
	"fmt"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s service) Routers(e *gin.Engine) {
	r := e.Group("/auth/google",
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
	callbackURL := fmt.Sprintf("%s://%s", schema, c.Request.Host+"/auth/google/callback")
	redirectURL, err := s.Login(c.Request.Context(), callbackURL)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	_ = loginResponse{
		RedirectURL: redirectURL,
		Err:         err,
	}
	if err != nil {
		e.EncodeError(c, err)
	}
	c.Redirect(301, redirectURL)
	// encodeResponse(c, sr)
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
	callbackURL := fmt.Sprintf("%s://%s", schema, c.Request.Host+"/auth/google/callback")
	tokenAccess, tokenRefresh, err := s.Callback(c.Request.Context(), callbackURL, dr)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	_ = callbackResponse{
		AccessToken:  tokenAccess.Token,
		RefreshToken: tokenRefresh.Token,
		ExpiresIn:    tokenRefresh.ExpiresIn,
		Err:          err,
	}

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

	c.Redirect(301, "http://localhost:4200")
}
