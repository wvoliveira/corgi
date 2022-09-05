package health

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func (s service) NewHTTP(r *mux.Router) {
	rr := r.PathPrefix("/health").Subrouter()
	// rr.Use(middleware.Authorizer(s.enforce))

	rr.HandleFunc("/ping", s.HTTPHealth).Methods("GET")
	//r.GET("/live", s.httpLive)
	//r.GET("/ready", s.httpReady)
}

func (s service) HTTPHealth(w http.ResponseWriter, r *http.Request) {
	response.Default(w, "pong", "", http.StatusOK)
}
