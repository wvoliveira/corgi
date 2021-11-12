package server

// The "account" is just over HTTP, so we just have a single transport.go.

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

/*
	MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
*/
func MakeHTTPHandler(s service, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	hAu := handlersAuth{s}
	hAc := handlersAccount{s}

	/*
		Auth with password endpoints and methods.
	*/
	r.HandleFunc("/api/v1/signin", hAu.signIn).Methods("POST")
	r.HandleFunc("/api/v1/signup", hAu.signUp).Methods("POST")

	/*
		Account endpoints and methods.
	*/
	r.Handle("/api/v1/accounts", addAccountHandler).Methods("POST")
	r.Handle("/api/v1/accounts/{id}", findAccountByIDHandler).Methods("GET")
	r.Handle("/api/v1/accounts", findAccountsHandler).Methods("GET")
	r.Handle("/api/v1/accounts/{id}", updateOrCreateAccountHandler).Methods("PUT")
	r.Handle("/api/v1/accounts/{id}", updateAccountHandler).Methods("PATCH")
	r.Handle("/api/v1/accounts/{id}", deleteAccountHandler).Methods("DELETE")

	/*
		URL endpoints and methods.
	*/
	r.Handle("/api/v1/urls", addURLHandler).Methods("POST")
	r.Handle("/api/v1/urls/{id}", findURLByIDHandler).Methods("GET")
	r.Handle("/api/v1/urls", findURLsHandler).Methods("GET")
	r.Handle("/api/v1/urls/{id}", updateOrCreateURLHandler).Methods("PUT")
	r.Handle("/api/v1/urls/{id}", updateURLHandler).Methods("PATCH")
	r.Handle("/api/v1/urls/{id}", deleteURLHandler).Methods("DELETE")

	return r
}

/*
	Auth handlers.
*/
type handlersAuth struct {
	s service
}

func (h handlersAuth) signIn(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sr)
}

func (h handlersAuth) signUp(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sr)
}

/*
	Account handlers.
*/
type handlersAccount struct {
	s service
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sr)
}

// FindAccountByID implements Service. Primarily useful in a client.
func (h handlersAccount) FindAccountByID(w http.ResponseWriter, r *http.Request) {
	request := findAccountByIDRequest{ID: id}
	response, err := e.FindAccountByIDEndpoint(ctx, request)
	if err != nil {
		return Account{}, err
	}
	resp := response.(findAccountByIDResponse)
	return resp.Account, resp.Err
}

// FindAccounts implements Service. Primarily useful in a client.
func (h handlersAccount) FindAccounts(w http.ResponseWriter, r *http.Request) {
	request := findAccountsRequest{Offset: offset, PageSize: pageSize}
	response, err := e.FindAccountsEndpoint(ctx, request)
	if err != nil {
		return []Account{}, err
	}
	resp := response.(findAccountsResponse)
	return resp.Accounts, resp.Err
}

// UpdateOrCreateAccount implements Service. Primarily useful in a client.
func (h handlersAccount) UpdateOrCreateAccount(w http.ResponseWriter, r *http.Request) {
	request := updateOrCreateAccountRequest{ID: id, Account: p}
	response, err := e.UpdateOrCreateAccountEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(updateOrCreateAccountResponse)
	return resp.Err
}

// UpdateAccount implements Service. Primarily useful in a client.
func (h handlersAccount) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	request := updateAccountRequest{ID: id, Account: p}
	response, err := e.UpdateAccountEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(updateAccountResponse)
	return resp.Err
}

// DeleteAccount implements Service. Primarily useful in a client.
func (h handlersAccount) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	request := deleteAccountRequest{ID: id}
	response, err := e.DeleteAccountEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(deleteAccountResponse)
	return resp.Err
}

/*
	URL enconde and decode.
*/

// AddURL implements Service. Primarily useful in a client.
func (e Endpoints) AddURL(ctx context.Context, p URL) error {
	request := addURLRequest{URL: p}
	response, err := e.AddURLEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(addURLResponse)
	return resp.Err
}

// FindURLByID implements Service. Primarily useful in a client.
func (e Endpoints) FindURLByID(ctx context.Context, id string) (URL, error) {
	request := findURLByIDRequest{ID: id}
	response, err := e.FindURLByIDEndpoint(ctx, request)
	if err != nil {
		return URL{}, err
	}
	resp := response.(findURLByIDResponse)
	return resp.URL, resp.Err
}

// FindURLs implements Service. Primarily useful in a client.
func (e Endpoints) FindURLs(ctx context.Context, offset, pageSize int) ([]URL, error) {
	request := findURLsRequest{Offset: offset, PageSize: pageSize}
	response, err := e.FindURLsEndpoint(ctx, request)
	if err != nil {
		return []URL{}, err
	}
	resp := response.(findURLsResponse)
	return resp.URLs, resp.Err
}

// UpdateOrCreateURL implements Service. Primarily useful in a client.
func (e Endpoints) UpdateOrCreateURL(ctx context.Context, id string, p URL) error {
	request := updateOrCreateURLRequest{ID: id, URL: p}
	response, err := e.UpdateOrCreateURLEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(updateOrCreateURLResponse)
	return resp.Err
}

// UpdateURL implements Service. Primarily useful in a client.
func (e Endpoints) UpdateURL(ctx context.Context, id string, p URL) error {
	request := updateURLRequest{ID: id, URL: p}
	response, err := e.UpdateURLEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(updateURLResponse)
	return resp.Err
}

// DeleteURL implements Service. Primarily useful in a client.
func (e Endpoints) DeleteURL(ctx context.Context, id string) error {
	request := deleteURLRequest{ID: id}
	response, err := e.DeleteURLEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(deleteURLResponse)
	return resp.Err
}

// MakeAddURLEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeAddURLEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addURLRequest)
		u, e := s.AddURL(ctx, req.URL)
		return addURLResponse{ID: u.ID, Keyword: u.Keyword, URL: u.URL, Title: u.Title, Err: e}, nil
	}
}

// MakeFindURLByIDEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeFindURLByIDEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(findURLByIDRequest)
		p, e := s.FindURLByID(ctx, req.ID)
		return findURLByIDResponse{URL: p, Err: e}, nil
	}
}

// MakeFindURLsEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeFindURLsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(findURLsRequest)
		p, e := s.FindURLs(ctx, req.Offset, req.PageSize)
		return findURLsResponse{URLs: p, Err: e}, nil
	}
}

// MakeUpdateOrCreateURLEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeUpdateOrCreateURLEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateOrCreateURLRequest)
		e := s.UpdateOrCreateURL(ctx, req.ID, req.URL)
		return updateOrCreateURLResponse{Err: e}, nil
	}
}

// MakeUpdateURLEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeUpdateURLEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateURLRequest)
		e := s.UpdateURL(ctx, req.ID, req.URL)
		return updateURLResponse{Err: e}, nil
	}
}

// MakeDeleteURLEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeDeleteURLEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteURLRequest)
		e := s.DeleteURL(ctx, req.ID)
		return deleteURLResponse{Err: e}, nil
	}
}
