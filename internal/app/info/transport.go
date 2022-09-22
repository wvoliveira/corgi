package info

import (
	"net/http"

	"github.com/gorilla/mux"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func (s service) NewHTTP(r *mux.Router) {
	rr := r.PathPrefix("/").Subrouter()
	//

	rr.HandleFunc("/info", s.HTTPInfo).Methods("GET")
	//r.GET("/live", s.httpLive)
	//r.GET("/ready", s.httpReady)
}

func (s service) HTTPInfo(w http.ResponseWriter, r *http.Request) {
	data, err := s.Info(r.Context())
	if err != nil {
		e.EncodeError(w, err)
		return
	}
	response.Default(w, data, "", http.StatusOK)
}
