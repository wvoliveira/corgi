package link

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/elga-io/corgi/internal/app/entity"
	"github.com/gorilla/mux"
)

type addRequest struct {
	Domain  string `json:"domain"`
	Keyword string `json:"keyword"`
	URL     string `json:"url"`
	Title   string `json:"title"`
	UserID  string `json:"user_id"`
}

type findByIDRequest struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type findAllRequest struct {
	Page   int    `json:"page"`
	Sort   string `json:"sort"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	UserID string `json:"user_id"`
}

type updateRequest struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Domain    string    `json:"domain"`
	Keyword   string    `json:"keyword"`
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	Active    string    `json:"active"`
	UserID    string    `json:"user_id"`
}

type deleteRequest struct {
	ID      string `json:"id"`
	Domain  string `json:"domain"`
	Keyword string `json:"keyword"`
	UserID  string `json:"user_id"`
}

func decodeAdd(r *http.Request) (req addRequest, err error) {
	ctx := r.Context()

	data := ctx.Value(entity.IdentityInfo{})
	if data == nil {
		return req, errors.New("impossible to get identity from context")
	}

	ii := data.(entity.IdentityInfo)

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, err
	}

	req.UserID = ii.UserID
	return req, nil
}

func decodeFindByID(r *http.Request) (req findByIDRequest, err error) {
	ctx := r.Context()
	vars := mux.Vars(r)

	data := ctx.Value(entity.IdentityInfo{})
	if data == nil {
		return req, errors.New("impossible to get identity from context")
	}

	ii := data.(entity.IdentityInfo)

	linkID := vars["id"]
	if linkID == "" {
		return req, errors.New("impossible to get link id")
	}

	req.ID = linkID
	req.UserID = ii.UserID
	return req, nil
}

func decodeFindAll(r *http.Request) (req findAllRequest, err error) {
	ctx := r.Context()
	params := r.URL.Query()

	data := ctx.Value(entity.IdentityInfo{})
	if data == nil {
		return req, errors.New("impossible to get identity from context")
	}

	ii := data.(entity.IdentityInfo)

	page, _ := strconv.Atoi(params.Get("page"))
	limit, _ := strconv.Atoi(params.Get("limit"))
	sort := params.Get("sort")

	if page == 0 {
		page = 1
	}
	if sort == "" {
		sort = "ID desc"
	}

	switch {
	case limit > 100:
		limit = 100
	case limit <= 0:
		limit = 10
	}
	offset := (page - 1) * limit

	req.Page = page
	req.Sort = sort
	req.Limit = limit
	req.Offset = offset
	req.UserID = ii.UserID
	return req, nil
}

func decodeUpdate(r *http.Request) (req updateRequest, err error) {
	ctx := r.Context()
	vars := mux.Vars(r)

	data := ctx.Value(entity.IdentityInfo{})
	if data == nil {
		return req, errors.New("impossible to get identity from context")
	}

	ii := data.(entity.IdentityInfo)

	linkID := vars["id"]
	if linkID == "" {
		return req, errors.New("impossible to get link id")
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, err
	}

	req.ID = linkID
	req.UserID = ii.UserID
	return req, nil
}

func decodeDelete(r *http.Request) (req deleteRequest, err error) {
	ctx := r.Context()
	vars := mux.Vars(r)

	data := ctx.Value(entity.IdentityInfo{})
	if data == nil {
		return req, errors.New("impossible to get identity from context")
	}

	ii := data.(entity.IdentityInfo)

	linkID := vars["id"]
	if linkID == "" {
		return req, errors.New("impossible to get link id")
	}

	req.ID = linkID
	req.UserID = ii.UserID
	return req, nil
}