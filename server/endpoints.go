package server

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

/*
	Endpoints collects all of the endpoints that compose a Endpoints service.
*/
type Endpoints struct{}

/*
	Auth sign-in and sign-up.
*/

// SignIn encode and decode.
func (e Endpoints) SignIn(ctx context.Context, a Account) error {
	// Create a request object.
	request := signInRequest{Email: a.Email, Password: a.Password, Session: a.Session}

	// Pass context and request to service.
	response, err := e.SignInEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(signInResponse)
	return resp.Err
}

// SignUP encode and decode.
func (e Endpoints) SignUP(ctx context.Context, a Account) error {
	request := signUpRequest{Email: a.Email, Password: a.Password}
	response, err := e.SignUpEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(signUpResponse)
	return resp.Err
}

// MakeSignInEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeSignInEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(signInRequest)
		p, e := s.SignIn(ctx, Account{Email: req.Email, Password: req.Password, Session: req.Session})
		return signInResponse{Session: p.Session, Err: e}, nil
	}
}

// MakeSignUpEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeSignUpEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(signUpRequest)
		e := s.SignUp(ctx, Account{Email: req.Email, Password: req.Password})
		return signUpResponse{Err: e}, nil
	}
}

/*
	Account encode and decode.
*/

// AddAccount implements Service. Primarily useful in a client.
func (e Endpoints) AddAccount(ctx context.Context, p Account) error {
	request := addAccountRequest{Account: p}
	response, err := e.AddAccountEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(addAccountResponse)
	return resp.Err
}

// FindAccountByID implements Service. Primarily useful in a client.
func (e Endpoints) FindAccountByID(ctx context.Context, id string) (Account, error) {
	request := findAccountByIDRequest{ID: id}
	response, err := e.FindAccountByIDEndpoint(ctx, request)
	if err != nil {
		return Account{}, err
	}
	resp := response.(findAccountByIDResponse)
	return resp.Account, resp.Err
}

// FindAccounts implements Service. Primarily useful in a client.
func (e Endpoints) FindAccounts(ctx context.Context, offset, pageSize int) ([]Account, error) {
	request := findAccountsRequest{Offset: offset, PageSize: pageSize}
	response, err := e.FindAccountsEndpoint(ctx, request)
	if err != nil {
		return []Account{}, err
	}
	resp := response.(findAccountsResponse)
	return resp.Accounts, resp.Err
}

// UpdateOrCreateAccount implements Service. Primarily useful in a client.
func (e Endpoints) UpdateOrCreateAccount(ctx context.Context, id string, p Account) error {
	request := updateOrCreateAccountRequest{ID: id, Account: p}
	response, err := e.UpdateOrCreateAccountEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(updateOrCreateAccountResponse)
	return resp.Err
}

// UpdateAccount implements Service. Primarily useful in a client.
func (e Endpoints) UpdateAccount(ctx context.Context, id string, p Account) error {
	request := updateAccountRequest{ID: id, Account: p}
	response, err := e.UpdateAccountEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(updateAccountResponse)
	return resp.Err
}

// DeleteAccount implements Service. Primarily useful in a client.
func (e Endpoints) DeleteAccount(ctx context.Context, id string) error {
	request := deleteAccountRequest{ID: id}
	response, err := e.DeleteAccountEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(deleteAccountResponse)
	return resp.Err
}

// MakeAddAccountEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeAddAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addAccountRequest)
		p, e := s.AddAccount(ctx, req.Account)
		return addAccountResponse{ID: p.ID, Email: p.Email, Err: e}, nil
	}
}

// MakeFindAccountByIDEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeFindAccountByIDEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(findAccountByIDRequest)
		p, e := s.FindAccountByID(ctx, req.ID)
		return findAccountByIDResponse{Account: p, Err: e}, nil
	}
}

// MakeFindAccountsEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeFindAccountsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(findAccountsRequest)
		p, e := s.FindAccounts(ctx, req.Offset, req.PageSize)
		return findAccountsResponse{Accounts: p, Err: e}, nil
	}
}

// MakeUpdateOrCreateAccountEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeUpdateOrCreateAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateOrCreateAccountRequest)
		e := s.UpdateOrCreateAccount(ctx, req.ID, req.Account)
		return updateOrCreateAccountResponse{Err: e}, nil
	}
}

// MakeUpdateAccountEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeUpdateAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateAccountRequest)
		e := s.UpdateAccount(ctx, req.ID, req.Account)
		return updateAccountResponse{Err: e}, nil
	}
}

// MakeDeleteAccountEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeDeleteAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteAccountRequest)
		e := s.DeleteAccount(ctx, req.ID)
		return deleteAccountResponse{Err: e}, nil
	}
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
