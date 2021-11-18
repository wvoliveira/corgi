package server

import (
	"encoding/json"
	"net/http"
)

func encodeResponse(w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		/*
			Not a Go kit transport error, but a business-logic error.
			Provide those as HTTP errors.
		*/
		encodeError(e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
