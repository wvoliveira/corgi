package group

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
)

type addRequest struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Description string   `json:"description"`
	UserIDs     []string `json:"user_Ids"`
}

type listRequest struct {
	Page   int    `json:"page"`
	Sort   string `json:"sort"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	UserID string `json:"-"`
}

type findByIDRequest struct {
	ID     string `json:"id"`
	UserID string `json:"-"`
}

func decodeAdd(r *http.Request) (payload addRequest, userID string, err error) {
	ctx := r.Context()

	data := ctx.Value(entity.IdentityInfo{})
	if data == nil {
		err = errors.New("impossible to get identity/user from context")
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return
	}

	payload.Name = strings.ToLower(payload.Name)
	err = checkName(payload.Name)
	if err != nil {
		return
	}

	identity := data.(entity.IdentityInfo)
	userID = identity.UserID
	return
}

func decodeList(r *http.Request) (request listRequest, err error) {
	ctx := r.Context()
	params := r.URL.Query()

	data := ctx.Value(entity.IdentityInfo{})
	if data == nil {
		err = e.ErrUserNotFoundInContext
		return
	}

	ii := data.(entity.IdentityInfo)

	page, _ := strconv.Atoi(params.Get("page"))
	limit, _ := strconv.Atoi(params.Get("limit"))
	sort := params.Get("sort")

	if page == 0 {
		page = 1
	}
	if sort == "" {
		sort = "name ASC"
	}

	switch {
	case limit > 100:
		limit = 100
	case limit <= 0:
		limit = 10
	}
	offset := (page - 1) * limit

	request.Page = page
	request.Sort = sort
	request.Limit = limit
	request.Offset = offset
	request.UserID = ii.UserID
	return
}

func decodeFindByID(r *http.Request) (req findByIDRequest, err error) {
	ctx := r.Context()
	vars := mux.Vars(r)

	data := ctx.Value(entity.IdentityInfo{})
	if data == nil {
		return req, errors.New("impossible to get identity from context")
	}

	identity := data.(entity.IdentityInfo)

	GroupID := vars["id"]
	if GroupID == "" {
		return req, errors.New("impossible to get group id")
	}

	req.ID = GroupID
	req.UserID = identity.UserID
	return req, nil
}
