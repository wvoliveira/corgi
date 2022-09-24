package group

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/wvoliveira/corgi/internal/pkg/entity"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
)

type addRequest struct {
	Name        string   `json:"name"`
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

func decodeAdd(r *http.Request) (request addRequest, userID string, err error) {
	ctx := r.Context()

	data := ctx.Value(entity.IdentityInfo{})
	if data == nil {
		err = errors.New("impossible to get identity/user from context")
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		return
	}

	// Group name is required. Without this, an error must happens.
	if request.Name == "" {
		err = errors.New("you must pass group name")
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
		sort = "ID desc"
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
