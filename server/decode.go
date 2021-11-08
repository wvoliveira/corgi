package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func decodeSignInRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req signInRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeSignUpRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req signUpRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeAddAccountRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req addAccountRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Account); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeFindAccountByIDRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return findAccountByIDRequest{ID: id}, nil
}

func decodeFindAccountsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page == 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return findAccountsRequest{Offset: offset, PageSize: pageSize}, nil
}

func decodeUpdateOrCreateAccountRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var account Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		return nil, err
	}
	return updateOrCreateAccountRequest{
		ID:      id,
		Account: account,
	}, nil
}

func decodeUpdateAccountRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var account Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		return nil, err
	}
	return updateAccountRequest{
		ID:      id,
		Account: account,
	}, nil
}

func decodeDeleteAccountRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return deleteAccountRequest{ID: id}, nil
}

func decodeAddURLRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req addURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req.URL); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeFindURLByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return findURLByIDRequest{ID: id}, nil
}

func decodeFindURLsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page == 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return findURLsRequest{Offset: offset, PageSize: pageSize}, nil
}

func decodeUpdateOrCreateURLRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var url URL
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		return nil, err
	}
	return updateOrCreateURLRequest{
		ID:  id,
		URL: url,
	}, nil
}

func decodeUpdateURLRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var url URL
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		return nil, err
	}
	return updateURLRequest{
		ID:  id,
		URL: url,
	}, nil
}

func decodeDeleteURLRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return deleteURLRequest{ID: id}, nil
}
