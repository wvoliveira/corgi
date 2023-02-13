package errors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// nolint
var (
	/**
		Internal errors.
	**/

	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
	ErrOnlyAdmin       = errors.New("only admin can do it")
	ErrEmailNotValid   = errors.New("try a valid e-mail")

	ErrFieldsRequired = errors.New("require more body fields for this request")

	// Token errors.
	ErrTokenInvalid = errors.New("invalid token")

	// ErrUnauthorized default authentication error.
	ErrUnauthorized     = errors.New("sorry, you are not unauthorized")
	ErrParseToken       = errors.New("there was an error in parsing token")
	ErrTokenExpired     = errors.New("your token has been expired")
	ErrNoTokenFound     = errors.New("token authorization not found in header")
	ErrAuthHeaderFormat = errors.New("must provide Authorization header with format `Bearer {token}`")

	/**
		User errors.
	**/

	// ErrUserNotFound error when user not found in database.
	ErrUserNotFound = errors.New("user not found")

	// ErrUserDeleteYourSelf user admin or user with permission with that cannot delete yourself.
	ErrUserDeleteYourSelf = errors.New("delete yourself? this is not a good idea")

	// ErrUserNotFoundInContext impossible to get identity or user from context of request.
	ErrUserNotFoundInContext = errors.New("impossible to get identity/user from context")

	// ErrUserFromSession when get user from session.
	ErrUserFromSession = errors.New("impossible to get user from session")

	/**
		Auth/password errors.
	**/

	// ErrAuthPasswordInternalError when an unkown error happens in auth/password category.
	ErrAuthPasswordInternalError = errors.New("unknown error happens when register user. Sorry about that")

	// ErrAuthPasswordUserAlreadyExists when anyone try to register with same e-mail.
	ErrAuthPasswordUserAlreadyExists = errors.New("this e-mail already exists in our database. Try another one")

	/**
		Link errors.
	**/

	// ErrLinkNotFound link not found in database.
	ErrLinkNotFound            = errors.New("domain and keyword combination not found")
	ErrLinkAlreadyExists       = errors.New("this link keyword already exists in our database")
	ErrLinkInvalidDomain       = errors.New("try to input a valid domain")
	ErrLinkInvalidKeyword      = errors.New("try to input a valid keyword between 6 and 15 chars")
	ErrLinkKeywordNotPermitted = errors.New("this keyword is not permitted")
	ErrLinkInvalidURL          = errors.New("try to input a valid destination (URL)")

	// With anonymous access, we can not create a shortener link with same URL.
	ErrAnonymousURLAlreadyExists = errors.New("with anonymous access, we can not create a shortener link with same URL")

	// Internal errors.
	ErrInternalServerError = errors.New("internal server error")
	ErrBadRouting          = errors.New("inconsistent mapping between route and handler (programmer error)")

	// ErrRequestNeedBody error if client not send a body payload.
	ErrRequestNeedBody = errors.New("methods POST and PATCH needs a body payload")

	/**
		Group errors.
	**/

	// ErrGroupAlreadyExists error when user try to create a group with a existent group name.
	ErrGroupAlreadyExists       = errors.New("group with this name already exists. Choose another one")
	ErrGroupNotFound            = errors.New("group with this ID was not found")
	ErrGroupInviteAlreadyExists = errors.New("this invite already exists. You need wait for response user")
)

type response struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

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

	resp := response{
		Status: "error",
		Data:   nil,
	}

	if err == nil {
		panic("encodeError with nil error")
	}

	resp.Message = err.Error()
	c.JSON(codeFrom(err), resp)
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound, ErrLinkNotFound, ErrGroupNotFound:
		return http.StatusNotFound

	case ErrRequestNeedBody, ErrInconsistentIDs,
		ErrLinkInvalidDomain, ErrLinkInvalidKeyword, ErrLinkKeywordNotPermitted, ErrLinkInvalidURL:
		return http.StatusBadRequest

	case ErrAlreadyExists, ErrLinkAlreadyExists, ErrAnonymousURLAlreadyExists, ErrAuthPasswordUserAlreadyExists:
		return http.StatusConflict

	case ErrUnauthorized, ErrNoTokenFound, ErrParseToken, ErrTokenExpired:
		return http.StatusUnauthorized

	default:
		return http.StatusInternalServerError
	}
}
