package server

import (
	"time"
)

/*
	Account represents a single struct for Account.
	ID should be globally unique.
*/
type Account struct {
	ID        string    `json:"id" gorm:"primaryKey;" example:"eed7df28-5a16-46f0-b5bf-c26071a42ade"`
	CreatedAt time.Time `json:"created_at,omitempty" example:"2021-10-18T00:45:07.818344164-03:00"`
	UpdatedAt time.Time `json:"updated_at,omitempty" example:"2021-10-18T00:49:06.160059334-03:00"`
	LastLogin time.Time `json:"last_login,omitempty" example:"2021-10-20T00:50:00.100059334-03:00"`

	Email    string `json:"email" example:"oliveira@live.it"`
	Password string `json:"-" example:"supersecret"`
	Tags     string `json:"tags"`

	Session string `json:"-" gorm:"-"`
}

/*
	URL represents a single struct for URL.
	ID should be globally unique.
*/
type URL struct {
	ID        string    `json:"id" gorm:"primaryKey;" example:"eed7df28-5a16-46f0-b5bf-c26071a42ade"`
	Keyword   string    `json:"keyword" gorm:"index" example:"google"`
	URL       string    `json:"url" example:"https://www.google.com"`
	Title     string    `json:"title" example:"Google Home"`
	Active    *bool     `json:"active" gorm:"type:bool;default:true" example:"false"`
	AccountID string    `json:"-"`
	Account   Account   `json:"-"`
	CreatedAt time.Time `json:"created_at" example:"2021-10-18T00:45:07.818344164-03:00"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-10-18T00:49:06.160059334-03:00"`
}

/*
	We have two options to return errors from the business logic.

	We could return the error via the endpoint itself. That makes certain things
	a little bit easier, like providing non-200 HTTP responses to the client. But
	Go kit assumes that endpoint errors are (or may be treated as)
	transport-domain errors. For example, an endpoint error will count against a
	circuit breaker error count.

	Therefore, it's often better to return service (business logic) errors in the
	response object. This means we have to do a bit more work in the HTTP
	response encoder to detect e.g. a not-found error and provide a proper HTTP
	status code. That work is done with the errorer interface, in transport.go.
	Response types that may contain business-logic errors implement that
	interface.
*/

/*
	Sign-in request and response structs.
*/

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Session  string `json:"-"`
}

type signInResponse struct {
	Session string `json:"-"`
	Err     error  `json:"err,omitempty"`
}

func (r signInResponse) error() error { return r.Err }

/*
	Sign-up request and response structs.
*/

type signUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signUpResponse struct {
	Err error `json:"err,omitempty"`
}

func (r signUpResponse) error() error { return r.Err }

/*
	Add URL resquest and response structs.
*/

type addURLRequest struct {
	URL URL
}

type addURLResponse struct {
	ID      string `json:"id" example:"eed7df28-5a16-46f0-b5bf-c26071a42ade"`
	Keyword string `json:"keyword" example:"google"`
	URL     string `json:"url" example:"https://www.google.com"`
	Title   string `json:"title" example:"Google Home"`
	Err     error  `json:"err,omitempty"`
}

func (r addURLResponse) error() error { return r.Err }

/*
	Find URL resquest and response structs.
*/

type findURLByIDRequest struct {
	ID string
}

type findURLByIDResponse struct {
	URL URL   `json:"data,omitempty"`
	Err error `json:"error,omitempty"`
}

func (r findURLByIDResponse) error() error { return r.Err }

/*
	Find URLs request and response structs.
*/

type findURLsRequest struct {
	Offset   int
	PageSize int
}

type findURLsResponse struct {
	URLs []URL `json:"data,omitempty"`
	Err  error `json:"error,omitempty"`
}

func (r findURLsResponse) error() error { return r.Err }

/*
	Update or Create URL resquest and response structs.
*/

type updateOrCreateURLRequest struct {
	ID  string
	URL URL
}

type updateOrCreateURLResponse struct {
	Err error `json:"err,omitempty"`
}

func (r updateOrCreateURLResponse) error() error { return nil }

/*
	Update URL resquest and response structs.
*/

type updateURLRequest struct {
	ID  string
	URL URL
}

type updateURLResponse struct {
	Err error `json:"err,omitempty"`
}

func (r updateURLResponse) error() error { return r.Err }

/*
	Delete URL request and response structs.
*/

type deleteURLRequest struct {
	ID string
}

type deleteURLResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteURLResponse) error() error { return r.Err }

/*
	Add Account request and response structs.
*/

type addAccountRequest struct {
	Account Account
}

type addAccountResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Err   error  `json:"err,omitempty"`
}

func (r addAccountResponse) error() error { return r.Err }

/*
	Find Account by ID request and response structs.
*/

type findAccountByIDRequest struct {
	ID string
}

type findAccountByIDResponse struct {
	Account Account `json:"data,omitempty"`
	Err     error   `json:"error,omitempty"`
}

func (r findAccountByIDResponse) error() error { return r.Err }

/*
	Find Accounts request and response structs.
*/

type findAccountsRequest struct {
	Offset   int
	PageSize int
}

type findAccountsResponse struct {
	Accounts []Account `json:"data,omitempty"`
	Err      error     `json:"error,omitempty"`
}

func (r findAccountsResponse) error() error { return r.Err }

/*
	Update or create Account request and response structs.
*/

type updateOrCreateAccountRequest struct {
	ID      string
	Account Account
}

type updateOrCreateAccountResponse struct {
	Err error `json:"err,omitempty"`
}

func (r updateOrCreateAccountResponse) error() error { return nil }

/*
	Update Account request and response structs.
*/

type updateAccountRequest struct {
	ID      string
	Account Account
}

type updateAccountResponse struct {
	Err error `json:"err,omitempty"`
}

func (r updateAccountResponse) error() error { return r.Err }

/*
	Delete Account request and response structs.
*/

type deleteAccountRequest struct {
	ID string
}

type deleteAccountResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteAccountResponse) error() error { return r.Err }
