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
func MakeHTTPHandler(s Service) http.Handler {
	r := mux.NewRouter()

	hau := handlersAuth{s}
	hac := handlersAccount{s}
	hur := handlersURL{s}

	/*
		Auth with password endpoints and methods.
	*/
	r.HandleFunc("/api/v1/signin", hau.SignIn).Methods("POST")
	r.HandleFunc("/api/v1/signup", hau.SignUp).Methods("POST")

	/*
		Account endpoints and methods.
	*/
	r.HandleFunc("/api/v1/accounts", IsAuthorized(hac.AddAccount)).Methods("POST")
	r.HandleFunc("/api/v1/accounts/{id}", IsAuthorized(hac.FindAccountByID)).Methods("GET")
	r.HandleFunc("/api/v1/accounts", IsAuthorized(hac.FindAccounts)).Methods("GET")
	r.HandleFunc("/api/v1/accounts/{id}", IsAuthorized(hac.UpdateOrCreateAccount)).Methods("PUT")
	r.HandleFunc("/api/v1/accounts/{id}", IsAuthorized(hac.UpdateAccount)).Methods("PATCH")
	r.HandleFunc("/api/v1/accounts/{id}", IsAuthorized(hac.DeleteAccount)).Methods("DELETE")

	/*
		URL endpoints and methods.
	*/
	r.HandleFunc("/api/v1/urls", IsAuthorized(hur.AddURL)).Methods("POST")
	r.HandleFunc("/api/v1/urls/{id}", IsAuthorized(hur.FindURLByID)).Methods("GET")
	r.HandleFunc("/api/v1/urls", IsAuthorized(hur.FindURLs)).Methods("GET")
	r.HandleFunc("/api/v1/urls/{id}", IsAuthorized(hur.UpdateOrCreateURL)).Methods("PUT")
	r.HandleFunc("/api/v1/urls/{id}", IsAuthorized(hur.UpdateURL)).Methods("PATCH")
	r.HandleFunc("/api/v1/urls/{id}", IsAuthorized(hur.DeleteURL)).Methods("DELETE")

	return r
}

/*
	Auth handlers.
*/
type handlersAuth struct {
	s Service
}

func (h handlersAuth) SignIn(w http.ResponseWriter, r *http.Request) {
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

func (h handlersAuth) SignUp(w http.ResponseWriter, r *http.Request) {
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
	de, err := decodeAddAccountRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	account, err := h.s.AddAccount(Account{Email: de.Email, Password: de.Password})
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
	de, err := decodeFindAccountByIDRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	account, err := h.s.FindAccountByID(de.ID)
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
	de, err := decodeFindAccountsRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	accounts, err := h.s.FindAccounts(de.Offset, de.PageSize)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Encode object to answer request (response).
	sr := findAccountsResponse{Accounts: accounts, Err: err}
	_ = encodeResponse(w, sr)
}

func (h handlersAccount) UpdateOrCreateAccount(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	de, err := decodeUpdateOrCreateAccountRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	err = h.s.UpdateOrCreateAccount(de.ID, de.Account)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Encode object to answer request (response).
	sr := updateOrCreateAccountResponse{Err: err}
	_ = encodeResponse(w, sr)
}

func (h handlersAccount) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	de, err := decodeUpdateAccountRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	err = h.s.UpdateAccount(de.ID, de.Account)
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
	de, err := decodeDeleteAccountRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	err = h.s.DeleteAccount(de.ID)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Encode object to answer request (response).
	sr := deleteAccountResponse{Err: err}
	_ = encodeResponse(w, sr)
}

/*
	URL handlers.
*/

type handlersURL struct {
	s Service
}

func (h handlersURL) AddURL(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	deAccount, deURL, err := decodeAddURLRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	url, err := h.s.AddURL(deAccount, URL{Keyword: deURL.Keyword, URL: deURL.URL, Title: deURL.Title})
	if err != nil {
		encodeError(err, w)
		return
	}

	// Encode object to answer request (response).
	sr := addURLResponse{Keyword: url.Keyword, URL: url.URL, Title: url.Title, Err: err}
	_ = encodeResponse(w, sr)
}

func (h handlersURL) FindURLByID(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	deAccount, deURL, err := decodeFindURLByIDRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	url, err := h.s.FindURLByID(deAccount, deURL.ID)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Encode object to answer request (response).
	sr := findURLByIDResponse{URL: url, Err: err}
	_ = encodeResponse(w, sr)
}

func (h handlersURL) FindURLs(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	deAccount, deURL, err := decodeFindURLsRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	urls, err := h.s.FindURLs(deAccount, deURL.Offset, deURL.PageSize)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Encode object to answer request (response).
	sr := findURLsResponse{URLs: urls, Err: err}
	_ = encodeResponse(w, sr)
}

func (h handlersURL) UpdateOrCreateURL(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	deAccount, deURL, err := decodeUpdateOrCreateURLRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	err = h.s.UpdateOrCreateURL(deAccount, deURL.ID, deURL.URL)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Encode object to answer request (response).
	sr := updateOrCreateURLResponse{Err: err}
	_ = encodeResponse(w, sr)
}

func (h handlersURL) UpdateURL(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	deAccount, deURL, err := decodeUpdateURLRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	err = h.s.UpdateURL(deAccount, deURL.ID, deURL.URL)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Encode object to answer request (response).
	sr := updateURLResponse{Err: err}
	_ = encodeResponse(w, sr)
}

func (h handlersURL) DeleteURL(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	deAccount, deURL, err := decodeDeleteURLRequest(r)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Business logic.
	err = h.s.DeleteURL(deAccount, deURL.ID)
	if err != nil {
		encodeError(err, w)
		return
	}

	// Encode object to answer request (response).
	sr := deleteURLResponse{Err: err}
	_ = encodeResponse(w, sr)
}
