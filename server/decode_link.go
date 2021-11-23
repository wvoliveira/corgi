package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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

func decodeUpdateLinkRequest(r *http.Request) (acc Account, req updateLinkRequest, err error) {
	acc = getAccountFromHeaders(r)
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
	acc = getAccountFromHeaders(r)
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return
	}
	req.ID = id
	return
}
