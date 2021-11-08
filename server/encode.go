package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

func encodeSignInResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		/*
			Not a Go kit transport error, but a business-logic error.
			Provide those as HTTP errors.
		*/
		encodeError(ctx, e.error(), w)
		return nil
	}

	var sessionToken string

	// Get session_token
	e, ok := response.(signInResponse)
	if ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	sessionToken = e.SessionToken

	/*
		Finally, we set the client cookie for "session_token" as the session token we just generated
		we also set an expiry time of 120 seconds, the same as the cache.
	*/
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(300 * time.Second),
	})

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
