package user

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/middleware"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func (s service) NewHTTP(r *mux.Router) {
	rr := r.PathPrefix("/v1/user").Subrouter()
	rr.Use(middleware.Checks)
	rr.Use(middleware.Auth(s.secret))

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
