package errors

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
	ErrOnlyAdmin       = errors.New("only admin can do it")
	ErrEmailNotValid   = errors.New("try a valid e-mail")

	ErrFieldsRequired = errors.New("require more body fields for this request")

	// Auth errors.
	ErrUnauthorized = errors.New("sorry, you are not unauthorized")
	ErrParseToken   = errors.New("there was an error in parsing token")
	ErrTokenExpired = errors.New("your token has been expired")
	ErrNoTokenFound = errors.New("token authorization not found in header")

	// Account errors.
	ErrAccountDeleteYourSelf = errors.New("delete yourself? this is not a good idea")

	// Link errors.
	ErrLinkIDNotFound           = errors.New("link id not found")
	ErrLinkKeywordNotFound      = errors.New("link keyword not found")
	ErrLinkKeywordAlreadyExists = errors.New("this link keyword already exists in our database")

	// Internal errors.
	ErrInternalServerError = errors.New("internal server error")
	ErrBadRouting          = errors.New("inconsistent mapping between route and handler (programmer error)")
)

/*
	errorer is implemented by all concrete response types that may contain
	errors. It allows us to change the HTTP response code without needing to
	trigger an endpoint (transport-level) error. For more information, read the
	big comment in endpoints.go.
*/
type Errorer interface {
	Error() error
}

// EncodeError generate a response for errors.
func EncodeError(err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound, ErrLinkIDNotFound, ErrLinkKeywordNotFound:
		return http.StatusNotFound
	case ErrLinkKeywordAlreadyExists, ErrAlreadyExists:
		return http.StatusForbidden
	case ErrInconsistentIDs, ErrAccountDeleteYourSelf:
		return http.StatusBadRequest
	case ErrUnauthorized, ErrNoTokenFound, ErrParseToken, ErrTokenExpired:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
