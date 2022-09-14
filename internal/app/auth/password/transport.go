package password

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/wvoliveira/corgi/internal/app/entity"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/middleware"
)

func (s service) NewHTTP(r *mux.Router) {
	rr := r.PathPrefix("/v1/auth/password").Subrouter()
	rr.Use(middleware.Checks)

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
	// Decode request and create a object with it.
	dr, err := decodeRegisterRequest(r)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		e.EncodeError(w, err)
		return
	}

	// Business logic.
	err = s.Register(r.Context(), entity.Identity{
		Provider: "email",
		UID:      dr.Email,
		Password: dr.Password,
	})
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	// Encode response to send to final-user.
	err = encodeRegister(w)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		e.EncodeError(w, err)
		return
	}
}
