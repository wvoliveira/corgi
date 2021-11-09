package server

import (
	"context"
	"encoding/json"
	"net/http"
)

/*
  Encode for awnser requests.
*/

func encodeSignInResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		/*
			Not a Go kit transport error, but a business-logic error.
			Provide those as HTTP errors.
		*/
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

/*
	encodeResponse is the common method to encode all response types to the
	client. I chose to do it this way because, since we're using JSON, there's no
	reason to provide anything more specific. It's certainly possible to
	specialize on a per-response (per-method) basis.
*/
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		/*
			Not a Go kit transport error, but a business-logic error.
			Provide those as HTTP errors.
		*/
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
