package facebook

import (
	"fmt"
	"net/http"

	e "github.com/elga-io/corgi/internal/pkg/errors"
	"github.com/elga-io/corgi/internal/pkg/middleware"
	"github.com/gorilla/mux"
)

func (s service) NewHTTP(r *mux.Router) {
	rr := r.PathPrefix("/auth/facebook").Subrouter()
	rr.Use(middleware.Checks)
	rr.Use(middleware.Authorizer(s.enforce))

	rr.HandleFunc("/login", s.HTTPLogin).Methods("GET")
	rr.HandleFunc("/callback", s.HTTPCallback).Methods("GET")
}

func (s service) HTTPLogin(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	_, err := decodeLoginRequest(r)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	// Business logic.
	schema := "http"
	if r.TLS != nil {
		schema = "https"
	}
	callbackURL := fmt.Sprintf("%s://%s", schema, r.Host+"/auth/facebook/callback")
	redirectURL, err := s.Login(r.Context(), callbackURL)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	// Encode object to answer request (response).
	if err != nil {
		e.EncodeError(w, err)
	}
	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
}

func (s service) HTTPCallback(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	dr, err := decodeCallbackRequest(r)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	// Business logic.
	schema := "http"
	if r.TLS != nil {
		schema = "https"
	}
	callbackURL := fmt.Sprintf("%s://%s", schema, r.Host+"/auth/facebook/callback")
	tokenAccess, tokenRefresh, err := s.Callback(r.Context(), callbackURL, dr)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	// Encode object to answer request (response).
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

	http.SetCookie(w, &cookieAccess)
	http.SetCookie(w, &cookieRefresh)

	http.Redirect(w, r, s.cfg.App.RedirectURL, http.StatusMovedPermanently)
}
