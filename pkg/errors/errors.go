package errors

import (
	"errors"
	"github.com/gin-gonic/gin"
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

	// Token errors.
	ErrTokenInvalid = errors.New("invalid token")

	// Auth errors.
	ErrUnauthorized = errors.New("sorry, you are not unauthorized")
	ErrParseToken   = errors.New("there was an error in parsing token")
	ErrTokenExpired = errors.New("your token has been expired")
	ErrNoTokenFound = errors.New("token authorization not found in header")

	// ErrUserNotFound error when user not found in database.
	ErrUserNotFound = errors.New("user not found")
	// ErrAccountDeleteYourSelf user admin or user with permission with that cannot delete yourself.
	ErrAccountDeleteYourSelf = errors.New("delete yourself? this is not a good idea")

	// ErrLinkNotFound link not found in database.
	ErrLinkNotFound       = errors.New("domain and keyword combination not found")
	ErrLinkAlreadyExists  = errors.New("this link keyword already exists in our database")
	ErrLinkInvalidDomain  = errors.New("try to input a valid domain")
	ErrLinkInvalidKeyword = errors.New("try to input a valid keyword between 6 and 15 chars")
	ErrLinkInvalidURL     = errors.New("try to input a valid destination (URL)")

	// Internal errors.
	ErrInternalServerError = errors.New("internal server error")
	ErrBadRouting          = errors.New("inconsistent mapping between route and handler (programmer error)")

	// ErrRequestNeedBody error if client not send a body payload.
	ErrRequestNeedBody = errors.New("methods POST and PATCH needs a body payload")
)

/*
	Errorer is implemented by all concrete response types that may contain
	errors. It allows us to change the HTTP response code without needing to
	trigger an endpoint (transport-level) error. For more information, read the
	big comment in endpoints.go.
*/
type Errorer interface {
	Error() error
}

// EncodeError generate a response for errors.
func EncodeError(c *gin.Context, err error) {
	if err == nil {
		panic("encodeError with nil error")
	}
	c.AbortWithStatusJSON(codeFrom(err), gin.H{"message": err.Error()})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound, ErrLinkNotFound:
		return http.StatusNotFound
	case ErrLinkAlreadyExists, ErrAlreadyExists, ErrLinkInvalidDomain, ErrLinkInvalidKeyword, ErrLinkInvalidURL:
		return http.StatusForbidden
	case ErrInconsistentIDs, ErrAccountDeleteYourSelf:
		return http.StatusBadRequest
	case ErrUnauthorized, ErrNoTokenFound, ErrParseToken, ErrTokenExpired:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
