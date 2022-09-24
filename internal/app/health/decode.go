package health

import (
	"net/http"

	"github.com/gorilla/mux"
)

func decodeHealthAuthProvider(r *http.Request) (provider string, err error) {
	vars := mux.Vars(r)
	provider = vars["provider"]
	return provider, nil
}
