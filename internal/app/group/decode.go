package group

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/model"
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

func decodeAdd(c *gin.Context) (payload addRequest, userID string, err error) {
	session := sessions.Default(c)
	v := session.Get("user")

	if v == nil {
		return payload, userID, errors.New("impossible to get user from session")
	}

	if err = c.ShouldBindJSON(&payload); err != nil {
		return payload, userID, err
	}

	payload.Name = strings.ToLower(payload.Name)
	err = checkName(payload.Name)
	if err != nil {
		return
	}

	r.UserID = v.(model.User).ID
	return
}

func decodeList(c *gin.Context) (request listRequest, err error) {
	ctx := r.Context()
	params := r.URL.Query()

	data := ctx.Value(model.IdentityInfo{})
	if data == nil {
		err = e.ErrUserNotFoundInContext
		return
	}

	ii := data.(model.IdentityInfo)

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

func decodeFindByID(c *gin.Context) (req findByIDRequest, err error) {
	ctx := r.Context()
	vars := mux.Vars(r)

	data := ctx.Value(model.IdentityInfo{})
	if data == nil {
		return req, errors.New("impossible to get identity from context")
	}

	identity := data.(model.IdentityInfo)

	GroupID := vars["id"]
	if GroupID == "" {
		return req, errors.New("impossible to get group id")
	}

	req.ID = GroupID
	req.UserID = identity.UserID
	return req, nil
}
