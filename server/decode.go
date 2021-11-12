package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/*
	Auth decodes.
*/

func decodeSignInRequest(r *http.Request) (req signInRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, err
	}
	return req, nil
}

func decodeSignUpRequest(r *http.Request) (req signUpRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, err
	}
	return req, nil
}

/*
	Account decodes.
*/

func decodeAddAccountRequest(r *http.Request) (req addAccountRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, err
	}
	return req, nil
}

func decodeFindAccountByIDRequest(r *http.Request) (req findAccountByIDRequest, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return req, ErrBadRouting
	}
	return findAccountByIDRequest{ID: id}, nil
}

func decodeFindAccountsRequest(r *http.Request) (req findAccountsRequest, err error) {
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

func decodeUpdateOrCreateAccountRequest(r *http.Request) (req updateOrCreateAccountRequest, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return req, ErrBadRouting
	}

	var account Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		return req, err
	}

	return updateOrCreateAccountRequest{
		ID:      id,
		Account: account,
	}, nil
}

func decodeUpdateAccountRequest(r *http.Request) (req updateAccountRequest, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return req, ErrBadRouting
	}
	var account Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		return req, err
	}
	return updateAccountRequest{
		ID:      id,
		Account: account,
	}, nil
}

func decodeDeleteAccountRequest(r *http.Request) (req deleteAccountRequest, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return req, ErrBadRouting
	}
	return deleteAccountRequest{ID: id}, nil
}

/*
	URL decodes.
*/

func decodeAddURLRequest(r *http.Request) (acc Account, req addURLRequest, err error) {
	acc = getAccountFromHeaders(r)
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return acc, req, err
	}
	return acc, req, nil
}

func decodeFindURLByIDRequest(r *http.Request) (acc Account, req findURLByIDRequest, err error) {
	acc = getAccountFromHeaders(r)
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		return acc, req, ErrBadRouting
	}
	return acc, findURLByIDRequest{ID: id}, nil
}

func decodeFindURLsRequest(r *http.Request) (acc Account, req findURLsRequest, err error) {
	acc = getAccountFromHeaders(r)
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

	req.PageSize = pageSize
	req.Offset = offset

	return
}

func decodeUpdateOrCreateURLRequest(r *http.Request) (acc Account, req updateOrCreateURLRequest, err error) {
	acc = getAccountFromHeaders(r)
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		return acc, req, ErrBadRouting
	}
	req.ID = id

	if err := json.NewDecoder(r.Body).Decode(&req.URL); err != nil {
		return acc, req, err
	}

	return
}

func decodeUpdateURLRequest(r *http.Request) (acc Account, req updateURLRequest, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return
	}
	req.ID = id

	if err = json.NewDecoder(r.Body).Decode(&req.URL); err != nil {
		return
	}
	return
}

func decodeDeleteURLRequest(r *http.Request) (acc Account, req deleteURLRequest, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return
	}
	req.ID = id
	return
}

func getAccountFromHeaders(r *http.Request) (a Account) {
	a.ID = r.Header.Get("AccountID")
	a.Email = r.Header.Get("AccountEmail")
	a.Role = r.Header.Get("AccountRole")
	return
}
