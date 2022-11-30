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
	rr.HandleFunc("/live", s.HTTPLive).Methods("GET")
	rr.HandleFunc("/ready", s.HTTPReady).Methods("GET")
}

func (s service) HTTPHealth(w http.ResponseWriter, r *http.Request) {
	healths, err := s.Health(r.Context())
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	httpStatusCode := http.StatusOK
	for _, item := range healths {
		if item.Required && item.Status != "OK" {
			httpStatusCode = http.StatusServiceUnavailable
		}
	}

	response.Default(w, healths, "", httpStatusCode)
}

func (s service) HTTPLive(w http.ResponseWriter, r *http.Request) {
	response.Default(w, "Live", "", http.StatusOK)
}

func (s service) HTTPReady(w http.ResponseWriter, r *http.Request) {
	response.Default(w, "Ready", "", http.StatusOK)
}
