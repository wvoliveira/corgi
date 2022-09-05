package password

import (
	"encoding/json"
	"net/http"

	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func encodeRegister(w http.ResponseWriter) (err error) {
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response.Response{
		Status: "successful",
	})
	return
}
