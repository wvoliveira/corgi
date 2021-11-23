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

	Role string `json:"role"`
	Tags string `json:"tags"`

	Active string `json:"active"`

	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`

	Tokens []Token
	Links  []Link
}

type Token struct {
	ID        string    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastUse   time.Time `json:"last_use"`

	AccessToken  string    `json:"-" gorm:"-"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    time.Time `json:"expires_in"`
	AccountID    string    `json:"account_id" gorm:"index"`
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

	Active string `json:"active"`

	AccountID string `json:"account_id" gorm:"index"`
}

/*
	Refresh request and response structs.
*/

type TokenRefreshRequest struct {
	ID           string `json:"-"`
	RefreshToken string `json:"refresh_token"`
}

type tokenRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Err          error  `json:"err,omitempty"`
}

func (r tokenRefreshResponse) error() error { return r.Err }

/*
	Sign-in request and response structs.
*/

type authLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Err          error  `json:"err,omitempty"`
}

func (r authLoginResponse) error() error { return r.Err }

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
