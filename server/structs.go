package server

import (
	"time"
)

/*
	Account represents a single struct for Account.
	ID should be globally unique.
*/
type Account struct {
	ID        string    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastLogin time.Time `json:"last_login"`

	Name     string `json:"name"`
	Email    string `json:"email" gorm:"index"`
	Password string `json:"password"`

	Role  string `json:"role"`
	Tags  string `json:"tags"`
	Token string `json:"token"`

	Active string `json:"active"`

	Links []Link
}

/*
	Link represents a single struct for Link.
	ID should be globally unique.
*/
type Link struct {
	ID        string    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Keyword     string `json:"keyword" gorm:"index"`
	Destination string `json:"destination"`
	Title       string `json:"title"`

	AccountID string `json:"account_id" gorm:"index"`

	Active string `json:"active"`
}

/*
	Sign-in request and response structs.
*/

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signInResponse struct {
	Token string `json:"token"`
	Err   error  `json:"err,omitempty"`
}

func (r signInResponse) error() error { return r.Err }

/*
	Sign-up request and response structs.
*/

type signUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signUpResponse struct {
	Err error `json:"err,omitempty"`
}

func (r signUpResponse) error() error { return r.Err }

/*
	Add Link resquest and response structs.
*/

type addLinkRequest struct {
	Keyword     string `json:"keyword"`
	Destination string `json:"destination"`
	Title       string `json:"title"`
}

type addLinkResponse struct {
	ID          string `json:"id"`
	Keyword     string `json:"keyword"`
	Destination string `json:"destination"`
	Title       string `json:"title"`
	Err         error  `json:"err,omitempty"`
}

func (r addLinkResponse) error() error { return r.Err }

/*
	Find Link resquest and response structs.
*/

type findLinkByIDRequest struct {
	ID string
}

type findLinkByIDResponse struct {
	Link Link  `json:"data,omitempty"`
	Err  error `json:"error,omitempty"`
}

func (r findLinkByIDResponse) error() error { return r.Err }

/*
	Find Links request and response structs.
*/

type findLinksRequest struct {
	Offset   int
	PageSize int
}

type findLinksResponse struct {
	Links []Link `json:"data,omitempty"`
	Err   error  `json:"error,omitempty"`
}

func (r findLinksResponse) error() error { return r.Err }

/*
	Update or Create Link resquest and response structs.
*/

type updateOrCreateLinkRequest struct {
	ID   string
	Link Link
}

type updateOrCreateLinkResponse struct {
	Err error `json:"err,omitempty"`
}

func (r updateOrCreateLinkResponse) error() error { return nil }

/*
	Update Link resquest and response structs.
*/

type updateLinkRequest struct {
	ID   string
	Link Link
}

type updateLinkResponse struct {
	Err error `json:"err,omitempty"`
}

func (r updateLinkResponse) error() error { return r.Err }

/*
	Delete Link request and response structs.
*/

type deleteLinkRequest struct {
	ID string
}

type deleteLinkResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteLinkResponse) error() error { return r.Err }

/*
	Add Account request and response structs.
*/

type addAccountRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
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
