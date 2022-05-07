package password

import (
	"net/http"

	"github.com/elga-io/corgi/internal/app/entity"
	e "github.com/elga-io/corgi/internal/pkg/errors"
	"github.com/elga-io/corgi/internal/pkg/response"
	"github.com/gorilla/mux"
)

func (s service) NewHTTP(r *mux.Router) {
	rr := r.PathPrefix("/auth/password").Subrouter()
	// middlewares.Checks())

	rr.HandleFunc("/login", s.HTTPLogin).Methods("POST")
	rr.HandleFunc("/register", s.HTTPRegister).Methods("POST")
}

func (s service) HTTPLogin(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	dr, err := decodeLoginRequest(r)
	if err != nil {
		e.EncodeError(w, err)
		return
	}
	identity := entity.Identity{Provider: "email", UID: dr.Email, Password: dr.Password}

	// Business logic.
	tokenAccess, tokenRefresh, err := s.Login(r.Context(), identity)
	if err != nil {
		e.EncodeError(w, err)
		return
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

	http.SetCookie(w, &cookieAccess)
	http.SetCookie(w, &cookieRefresh)

	w.WriteHeader(200)
}

func (s service) HTTPRegister(w http.ResponseWriter, r *http.Request) {
	dr, err := decodeRegisterRequest(r)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	err = s.Register(r.Context(), entity.Identity{
		Provider: "email",
		UID:      dr.Email,
		Password: dr.Password,
	})
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	sr := registerResponse{Err: err}
	response.Default(w, sr, "", http.StatusNoContent)
}
