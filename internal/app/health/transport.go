package health

import (
	"net/http"

	"github.com/elga-io/corgi/internal/pkg/response"
	"github.com/gorilla/mux"
)

func (s service) NewHTTP(r *mux.Router) {
	rr := r.PathPrefix("/health").Subrouter()
	// middlewares.Authorizer(s.enforce))

	rr.HandleFunc("/ping", s.HTTPHealth).Methods("GET")
	//r.GET("/live", s.httpLive)
	//r.GET("/ready", s.httpReady)
}

func (s service) HTTPHealth(w http.ResponseWriter, r *http.Request) {
	response.Default(w, "pong "+s.version, "", http.StatusOK)
}
