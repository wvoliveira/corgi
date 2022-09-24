package health

import (
	"net/http"

	"github.com/gorilla/mux"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func (s service) NewHTTP(r *mux.Router) {
	rr := r.PathPrefix("/health").Subrouter()

	rr.HandleFunc("", s.HTTPHealth).Methods("GET")

	rr.HandleFunc("/database", s.HTTPHealthDatabase).Methods("GET")

	rr.HandleFunc("/auth", s.HTTPHealthAuth).Methods("GET")
	rr.HandleFunc("/auth/{provider}", s.HTTPHealthAuthProvider).Methods("GET")

	rr.HandleFunc("/live", s.HTTPHealthLive).Methods("GET")
	rr.HandleFunc("/ready", s.HTTPHealthReady).Methods("GET")
}

func (s service) HTTPHealth(w http.ResponseWriter, r *http.Request) {
	healths, err := s.Health(r.Context())
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	response.Default(w, healths, "", http.StatusOK)
}

func (s service) HTTPHealthDatabase(w http.ResponseWriter, r *http.Request) {
	health := s.HealthDatabase(r.Context())
	response.Default(w, health, "", http.StatusOK)
}

func (s service) HTTPHealthAuth(w http.ResponseWriter, r *http.Request) {
	healths, err := s.HealthAuth(r.Context(), nil)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	response.Default(w, healths, "", http.StatusOK)
}

func (s service) HTTPHealthAuthProvider(w http.ResponseWriter, r *http.Request) {
	provider, err := decodeHealthAuthProvider(r)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	healths, err := s.HealthAuth(r.Context(), []string{provider})
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	response.Default(w, healths[0], "", http.StatusOK)
}

func (s service) HTTPHealthLive(w http.ResponseWriter, r *http.Request) {
	response.Default(w, "Live", "", http.StatusOK)
}

func (s service) HTTPHealthReady(w http.ResponseWriter, r *http.Request) {
	s.HTTPHealth(w, r)
}
