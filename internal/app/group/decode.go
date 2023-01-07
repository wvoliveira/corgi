package group

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/pkg/util"
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
}

type findByIDRequest struct {
	ID string `json:"id"`
}

func decodeAdd(c *gin.Context) (payload addRequest, userID string, err error) {
	userID, err = util.GetUserFromSession(c)

	if err != nil {
		return
	}

	if err = c.ShouldBindJSON(&payload); err != nil {
		return payload, userID, err
	}

	payload.Name = strings.ToLower(payload.Name)
	err = checkName(payload.Name)
	if err != nil {
		return
	}
	return
}

func decodeList(c *gin.Context) (request listRequest, userID string, err error) {
	userID, err = util.GetUserFromSession(c)

	if err != nil {
		return
	}

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	sort := c.Query("sort")

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
	return
}

func decodeFindByID(c *gin.Context) (req findByIDRequest, userID string, err error) {
	userID, err = util.GetUserFromSession(c)

	if err != nil {
		return
	}

	GroupID := c.Param("id")

	if GroupID == "" {
		return req, userID, errors.New("impossible to get group id")
	}

	req.ID = GroupID
	return req, userID, err
}
