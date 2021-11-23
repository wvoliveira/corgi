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

func decodeAddAccountRequest(r *http.Request) (acc Account, req addAccountRequest, err error) {
	acc = getAccountFromHeaders(r)
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}
	return
}

func decodeFindAccountByIDRequest(r *http.Request) (acc Account, req findAccountByIDRequest, err error) {
	acc = getAccountFromHeaders(r)
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return acc, req, ErrBadRouting
	}
	req.ID = id
	return acc, req, nil
}

func decodeFindAccountsRequest(r *http.Request) (acc Account, req findAccountsRequest, err error) {
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

	req.Offset = offset
	req.PageSize = offset
	return acc, req, nil
}

func decodeUpdateOrCreateAccountRequest(r *http.Request) (acc Account, req updateOrCreateAccountRequest, err error) {
	acc = getAccountFromHeaders(r)
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return acc, req, ErrBadRouting
	}

	var accountPayload Account
	if err := json.NewDecoder(r.Body).Decode(&accountPayload); err != nil {
		return acc, req, err
	}
	req.ID = id
	req.Account = accountPayload
	return acc, req, nil
}

func decodeUpdateAccountRequest(r *http.Request) (acc Account, req updateAccountRequest, err error) {
	acc = getAccountFromHeaders(r)
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return acc, req, ErrBadRouting
	}
	var accountPayload Account
	if err := json.NewDecoder(r.Body).Decode(&accountPayload); err != nil {
		return acc, req, err
	}
	req.ID = id
	req.Account = accountPayload
	return acc, req, nil
}

func decodeDeleteAccountRequest(r *http.Request) (acc Account, req deleteAccountRequest, err error) {
	acc = getAccountFromHeaders(r)
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return acc, req, ErrBadRouting
	}
	req.ID = id
	return acc, req, nil
}

/*
	Link decodes.
*/

func decodeAddLinkRequest(r *http.Request) (acc Account, req addLinkRequest, err error) {
	acc = getAccountFromHeaders(r)
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return acc, req, err
	}
	return acc, req, nil
}

func decodeFindLinkByIDRequest(r *http.Request) (acc Account, req findLinkByIDRequest, err error) {
	acc = getAccountFromHeaders(r)
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		return acc, req, ErrBadRouting
	}
	return acc, findLinkByIDRequest{ID: id}, nil
}

func decodeFindLinksRequest(r *http.Request) (acc Account, req findLinksRequest, err error) {
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

func decodeUpdateOrCreateLinkRequest(r *http.Request) (acc Account, req updateOrCreateLinkRequest, err error) {
	acc = getAccountFromHeaders(r)
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		return acc, req, ErrBadRouting
	}
	req.ID = id

	if err := json.NewDecoder(r.Body).Decode(&req.Link); err != nil {
		return acc, req, err
	}

	return
}

func decodeUpdateLinkRequest(r *http.Request) (acc Account, req updateLinkRequest, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return
	}
	req.ID = id

	if err = json.NewDecoder(r.Body).Decode(&req.Link); err != nil {
		return
	}
	return
}

func decodeDeleteLinkRequest(r *http.Request) (acc Account, req deleteLinkRequest, err error) {
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
