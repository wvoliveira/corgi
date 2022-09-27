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

	if data == nil {
		return
	}

	statusText := "successful"
	if status >= 500 && status <= 599 {
		statusText = "error"
	}

	_ = json.NewEncoder(w).Encode(Response{
		Status:  statusText,
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
