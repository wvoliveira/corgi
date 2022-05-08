package auth

import (
	"net/http"

	e "github.com/elga-io/corgi/internal/pkg/errors"
	"github.com/elga-io/corgi/internal/pkg/response"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func (s service) NewHTTP(r *mux.Router) {
	rr := r.PathPrefix("/auth").Subrouter()
	// middlewares.Checks(),
	// middlewares.Auth(s.secret),
	// middlewares.Authorizer(s.enforce))

	rr.HandleFunc("/logout", s.HTTPLogout).Methods("GET")
}

func (s service) HTTPLogout(w http.ResponseWriter, r *http.Request) {
	l := log.Ctx(r.Context())

	// Decode request to object.
	dr, err := decodeLogout(r)
	if err != nil {
		l.Error().Caller().Msg(err.Error())
		e.EncodeError(w, err)
		return
	}

	// Business logic.
	err = s.Logout(r.Context(), dr.Token)
	if err != nil {
		l.Error().Caller().Msg(err.Error())
		e.EncodeError(w, err)
		return
	}

	cookieAccess := http.Cookie{Name: "access_token", MaxAge: -1}
	cookieRefresh := http.Cookie{Name: "refresh_token_id", MaxAge: -1}
	cookieLogged := http.Cookie{Name: "logged", MaxAge: -1}

	http.SetCookie(w, &cookieAccess)
	http.SetCookie(w, &cookieRefresh)
	http.SetCookie(w, &cookieLogged)

	// Encode object to answer request (response).
	sr := logoutResponse{Err: err}
	response.Default(w, sr, "", http.StatusNoContent)
}
