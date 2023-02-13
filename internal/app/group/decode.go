package group

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/pkg/common"
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

type deleteRequest struct {
	UserID  string `json:"-"`
	GroupID string `json:"-"`
}

type inviteAddRequest struct {
	GroupID   string `json:"group_id"`
	UserID    string `json:"user_id"`
	InvitedBy string `json:"-"`
}

func decodeAdd(c *gin.Context) (req addRequest, userID string, err error) {
	userID, err = common.GetUserFromSession(c)

	if err != nil {
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		return req, userID, err
	}

	req.Name = strings.ToLower(req.Name)
	err = checkName(req.Name)
	if err != nil {
		return
	}
	return
}

func decodeList(c *gin.Context) (req listRequest, userID string, err error) {
	userID, err = common.GetUserFromSession(c)

	if err != nil {
		return
	}

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	sort := c.Query("sort")

	if page == 0 {
		page = 1
	}

	// TODO: rule for "sort" content
	// like ASC or DESC
	if sort == "" {
		sort = "ASC"
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
	return
}

func decodeFindByID(c *gin.Context) (req findByIDRequest, userID string, err error) {
	userID, err = common.GetUserFromSession(c)

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

func decodeDelete(c *gin.Context) (req deleteRequest, err error) {
	userID, err := common.GetUserFromSession(c)
	if err != nil {
		return
	}

	GroupID := c.Param("id")
	if GroupID == "" {
		return req, errors.New("impossible to get group id from path")
	}

	req.UserID = userID
	req.GroupID = GroupID
	return
}

func decodeInviteAdd(c *gin.Context) (req inviteAddRequest, err error) {
	userID, err := common.GetUserFromSession(c)
	if err != nil {
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		return req, err
	}

	req.InvitedBy = userID

	if req.GroupID == "" || req.UserID == "" {
		err = errors.New("you need specific group_id and user_id to send invite to the group")
	}

	return
}
