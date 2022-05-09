package response

import (
	"encoding/json"
	"net/http"
)

// Response default response for http requests.
type Response struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func Default(w http.ResponseWriter, data interface{}, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Response{
		Status:  "successful",
		Data:    data,
		Message: message,
	})
}

func NotImplemented(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(Response{
		Status:  "successful",
		Message: "not implemented yet",
	})
}