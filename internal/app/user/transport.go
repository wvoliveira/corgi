package user

import (
	"net/http"

	"github.com/elga-io/corgi/internal/app/entity"
	e "github.com/elga-io/corgi/internal/pkg/errors"
	"github.com/elga-io/corgi/internal/pkg/response"
	"github.com/gorilla/mux"
)

func (s service) NewHTTP(r *mux.Router) {
	rr := r.PathPrefix("/api/v1/user").Subrouter()
	// middlewares.Checks(),
	// middlewares.Auth(s.secret),
	// middlewares.Authorizer(s.enforce))

	rr.HandleFunc("/me", s.HTTPFind).Methods("GET")
	rr.HandleFunc("/me", s.HTTPUpdate).Methods("PATCH")
}

func (s service) HTTPFind(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	dr, err := decodeFind(r)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	// Business logic.
	user, err := s.Find(r.Context(), dr.UserID)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	// Encode object to answer request (response).
	var identities []identity
	for _, i := range user.Identities {
		idt := identity{
			Provider: i.Provider,
			UID:      i.UID,
		}
		identities = append(identities, idt)
	}

	ur := userResponse{Name: user.Name, Role: user.Role, Identities: identities}
	response.Default(w, ur, "", http.StatusOK)
}

func (s service) HTTPUpdate(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	dr, err := decodeUpdate(r)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	// Business logic.
	user, err := s.Update(r.Context(), entity.User{ID: dr.UserID, Name: dr.Name})
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	// Encode object to answer request (response).
	sr := updateResponse{
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Err:       err,
	}
	response.Default(w, sr, "", http.StatusOK)
}
