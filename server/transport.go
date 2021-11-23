package server

// The "account" is just over HTTP, so we just have a single transport.go.

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

/*
	MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
*/
func MakeHTTPHandler(s Service, m Middlewares) http.Handler {
	r := mux.NewRouter()

	handlerAuth := handlersAuth{s}
	handlerAccount := handlersAccount{s}
	handlerLink := handlersLink{s}

	r.Use(m.Logging)

	// Auth subrouter.
	subAuth := r.PathPrefix("/api/v1/auth").Subrouter().StrictSlash(true)
	subAuth.HandleFunc("/password/login", handlerAuth.Login).Methods("POST")
	subAuth.HandleFunc("/password/register", handlerAuth.Register).Methods("POST")

	// Accounts subrouter.
	subAccounts := r.PathPrefix("/api/v1/accounts").Subrouter().StrictSlash(true)
	subAccounts.HandleFunc("/", handlerAccount.AddAccount).Methods("POST")
	subAccounts.HandleFunc("/{id}", handlerAccount.FindAccountByID).Methods("GET")
	subAccounts.HandleFunc("/", handlerAccount.FindAccounts).Methods("GET")
	subAccounts.HandleFunc("/{id}", handlerAccount.UpdateAccount).Methods("PATCH")
	subAccounts.HandleFunc("/{id}", handlerAccount.DeleteAccount).Methods("DELETE")
	subAccounts.Use(m.Authentication)

	// Links subrouter.
	subLinks := r.PathPrefix("/api/v1/links").Subrouter().StrictSlash(true)
	subLinks.HandleFunc("/", handlerLink.AddLink).Methods("POST")
	subLinks.HandleFunc("/{id}", handlerLink.FindLinkByID).Methods("GET")
	subLinks.HandleFunc("/", handlerLink.FindLinks).Methods("GET")
	subLinks.HandleFunc("/{id}", handlerLink.UpdateLink).Methods("PATCH")
	subLinks.HandleFunc("/{id}", handlerLink.DeleteLink).Methods("DELETE")
	subLinks.Use(m.Authentication)

	// Group all subrouters.
	r.Handle("/api/v1/auth", subAuth)
	r.Handle("/api/v1/accounts", subAccounts)
	r.Handle("/api/v1/links", subLinks)
	return r
}

/*
	Auth handlers.
*/
type handlersAuth struct {
	s Service
}

func (h handlersAuth) Login(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	dr, err := decodeSignInRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	account, err := h.s.SignIn(Account{Email: dr.Email, Password: dr.Password})
	if err != nil {
		encodeError(err, w)
		return
	}

	// Encode object to answer request (response).
	sr := signInResponse{Token: account.Token, Err: err}
	_ = encodeResponse(w, sr)
}

func (h handlersAuth) Register(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	dr, err := decodeSignUpRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	err = h.s.SignUp(Account{Name: dr.Name, Email: dr.Email, Password: dr.Password})
	if err != nil {
		encodeError(err, w)
		return
	}

	// Encode object to answer request (response).
	sr := signUpResponse{Err: err}
	_ = encodeResponse(w, sr)
}

/*
	Account handlers.
*/

type handlersAccount struct {
	s Service
}

func (h handlersAccount) AddAccount(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	auth, payload, err := decodeAddAccountRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	account, err := h.s.AddAccount(auth, Account{Email: payload.Email, Password: payload.Password})
	if err != nil {
		encodeError(err, w)
		return
	}

	// Encode object to answer request (response).
	sr := addAccountResponse{ID: account.ID, Email: account.Email, Err: err}
	_ = encodeResponse(w, sr)
}

func (h handlersAccount) FindAccountByID(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	auth, payload, err := decodeFindAccountByIDRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	account, err := h.s.FindAccountByID(auth, payload.ID)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Encode object to answer request (response).
	sr := findAccountByIDResponse{Account: account, Err: err}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sr)
}

func (h handlersAccount) FindAccounts(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	auth, payload, err := decodeFindAccountsRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	accounts, err := h.s.FindAccounts(auth, payload.Offset, payload.PageSize)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Encode object to answer request (response).
	sr := findAccountsResponse{Accounts: accounts, Err: err}
	_ = encodeResponse(w, sr)
}

func (h handlersAccount) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	auth, payload, err := decodeUpdateAccountRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	err = h.s.UpdateAccount(auth, payload.ID, payload.Account)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Encode object to answer request (response).
	sr := updateAccountResponse{Err: err}
	_ = encodeResponse(w, sr)
}

func (h handlersAccount) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	auth, payload, err := decodeDeleteAccountRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	err = h.s.DeleteAccount(auth, payload.ID)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Encode object to answer request (response).
	sr := deleteAccountResponse{Err: err}
	_ = encodeResponse(w, sr)
}

/*
	Link handlers.
*/

type handlersLink struct {
	s Service
}

func (h handlersLink) AddLink(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	decodeAccount, decodeLink, err := decodeAddLinkRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	link, err := h.s.AddLink(decodeAccount, Link{Keyword: decodeLink.Keyword, Destination: decodeLink.Destination, Title: decodeLink.Title})
	if err != nil {
		encodeError(err, w)
		return
	}

	// Encode object to answer request (response).
	sr := addLinkResponse{ID: link.ID, Keyword: link.Keyword, Destination: link.Destination, Title: link.Title, Err: err}
	_ = encodeResponse(w, sr)
}

func (h handlersLink) FindLinkByID(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	deAccount, deLink, err := decodeFindLinkByIDRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	link, err := h.s.FindLinkByID(deAccount, deLink.ID)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Encode object to answer request (response).
	sr := findLinkByIDResponse{Link: link, Err: err}
	_ = encodeResponse(w, sr)
}

func (h handlersLink) FindLinks(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	deAccount, deLink, err := decodeFindLinksRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	links, err := h.s.FindLinks(deAccount, deLink.Offset, deLink.PageSize)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Encode object to answer request (response).
	sr := findLinksResponse{Links: links, Err: err}
	_ = encodeResponse(w, sr)
}

func (h handlersLink) UpdateLink(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	deAccount, deLink, err := decodeUpdateLinkRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	err = h.s.UpdateLink(deAccount, deLink.ID, deLink.Link)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Encode object to answer request (response).
	sr := updateLinkResponse{Err: err}
	_ = encodeResponse(w, sr)
}

func (h handlersLink) DeleteLink(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	deAccount, deLink, err := decodeDeleteLinkRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	err = h.s.DeleteLink(deAccount, deLink.ID)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Encode object to answer request (response).
	sr := deleteLinkResponse{Err: err}
	_ = encodeResponse(w, sr)
}
