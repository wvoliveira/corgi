package server

import (
	"encoding/json"
	"errors"
	"net/http"
)

//nolint
var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")

	ErrFieldsRequired      = errors.New("fields required: email and password")
	ErrInternalServerError = errors.New("internal server error")
	ErrUnauthorized        = errors.New("unauthorized")

	ErrNoTokenFound = errors.New("no token found")
	ErrParseToken   = errors.New("there was an error in parsing token")
	ErrTokenExpired = errors.New("your token has been expired")

	/*
		ErrBadRouting is returned when an expected path variable is missing.
		It always indicates programmer error.
	*/
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

/*
	errorer is implemented by all concrete response types that may contain
	errors. It allows us to change the HTTP response code without needing to
	trigger an endpoint (transport-level) error. For more information, read the
	big comment in endpoints.go.
*/
type errorer interface {
	error() error
}

func encodeError(err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists, ErrInconsistentIDs:
		return http.StatusBadRequest
	case ErrUnauthorized, ErrNoTokenFound, ErrParseToken, ErrTokenExpired:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
