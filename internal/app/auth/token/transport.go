package token

import (
	"net/http"

	"github.com/elga-io/corgi/internal/app/entity"
	e "github.com/elga-io/corgi/internal/pkg/errors"
	"github.com/elga-io/corgi/internal/pkg/response"
	"github.com/gorilla/mux"
)

func (s service) NewHTTP(r *mux.Router) {
	rr := r.PathPrefix("/auth/token").Subrouter()
	// middlewares.Checks(),
	// middlewares.Auth(s.secret),
	// middlewares.Authorizer(s.enforce))

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
