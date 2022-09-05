package token

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wvoliveira/corgi/internal/app/entity"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/middleware"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func (s service) NewHTTP(r *mux.Router) {
	rr := r.PathPrefix("/auth/token").Subrouter()
	rr.Use(middleware.Checks)
	rr.Use(middleware.Auth(s.secret))
	rr.Use(middleware.Authorizer(s.enforce))

	rr.HandleFunc("/refresh", s.HTTPRefresh).Methods("POST")
}

func (s service) HTTPRefresh(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	dr, err := decodeRefreshRequest(r)
	if err != nil {
		e.EncodeError(w, err)
		return
	}
	token := entity.Token{ID: dr.RefreshTokenID}

	// Business logic.
	tokenAccess, tokenRefresh, err := s.Refresh(r.Context(), token)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	// Encode object to answer request (response).
	sr := refreshResponse{
		AccessToken:  tokenAccess.Token,
		RefreshToken: tokenRefresh.Token,
		ExpiresIn:    tokenAccess.ExpiresIn,
		Err:          err,
	}
	// encodeResponse(w, sr)
	response.Default(w, sr, "", http.StatusNoContent)
}
