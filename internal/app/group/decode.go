package group

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

type addRequest struct {
	WhoID       string
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Description string   `json:"description"`
	UserIDs     []string `json:"user_Ids"`
}

type listRequest struct {
	WhoID  string
	Page   int    `json:"page"`
	Sort   string `json:"sort"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

type findByIDRequest struct {
	WhoID   string
	GroupID string `uri:"id" binding:"required"`
}

type deleteRequest struct {
	WhoID   string
	GroupID string `uri:"id" binding:"required"`
}

type invitesAddByIDRequest struct {
	WhoID     string
	InvitedBy string
	GroupID   string `uri:"id"`
	UserEmail string `json:"user_email" binding:"required"`
}

type invitesListByIDRequest struct {
	WhoID   string
	GroupID string `uri:"id"`
	Page    int    `form:"page"`
	Sort    string `form:"sort"`
	Offset  int    `form:"offset"`
	Limit   int    `form:"limit"`
}

type invitesListRequest struct {
	WhoID  string
	Page   int    `form:"page"`
	Sort   string `form:"sort"`
	Offset int    `form:"offset"`
	Limit  int    `form:"limit"`
}

func decodeAdd(c *gin.Context) (req addRequest, err error) {
	v, ok := c.Get("user_id")
	if !ok {
		err = errors.New("impossible to know who you are")
		return
	}

	req.WhoID = v.(string)
	err = c.ShouldBindJSON(&req)
	return
}

func decodeList(c *gin.Context) (req listRequest, err error) {
	v, ok := c.Get("user_id")
	if !ok {
		err = errors.New("impossible to know who you are")
		return
	}

	req.WhoID = v.(string)
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

func decodeFindByID(c *gin.Context) (req findByIDRequest, err error) {
	v, ok := c.Get("user_id")
	if !ok {
		err = errors.New("impossible to know who you are")
		return
	}

	err = c.ShouldBindUri(&req)
	if err != nil {
		return
	}

	req.WhoID = v.(string)
	return
}

func decodeDelete(c *gin.Context) (req deleteRequest, err error) {
	v, ok := c.Get("user_id")
	if !ok {
		err = errors.New("impossible to know who you are")
		return
	}

	if err = c.ShouldBindUri(&req); err != nil {
		return req, err
	}

	req.WhoID = v.(string)
	return
}

func decodeInvitesAddByID(c *gin.Context) (req invitesAddByIDRequest, err error) {
	v, ok := c.Get("user_id")
	if !ok {
		err = errors.New("impossible to know who you are")
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		return req, err
	}

	req.GroupID = c.Param("id")
	req.WhoID = v.(string)
	req.InvitedBy = req.WhoID
	return
}

func decodeInvitesListByID(c *gin.Context) (req invitesListByIDRequest, err error) {
	v, ok := c.Get("user_id")
	if !ok {
		err = errors.New("impossible to know who you are")
		return
	}

	req.WhoID = v.(string)
	if err = c.ShouldBindQuery(&req); err != nil {
		return req, err
	}

	if req.Page == 0 {
		req.Page = 1
	}

	// TODO: rule for "sort" content like ASC or DESC
	if req.Sort == "" {
		req.Sort = "ASC"
	}

	switch {
	case req.Limit > 100:
		req.Limit = 100
	case req.Limit <= 0:
		req.Limit = 10
	}
	offset := (req.Page - 1) * req.Limit

	req.GroupID = c.Param("id")
	req.Offset = offset
	return
}

func decodeInvitesList(c *gin.Context) (req invitesListRequest, err error) {
	v, ok := c.Get("user_id")
	if !ok {
		err = errors.New("impossible to know who you are")
		return
	}

	req.WhoID = v.(string)
	if err = c.ShouldBindQuery(&req); err != nil {
		return req, err
	}

	if req.Page == 0 {
		req.Page = 1
	}

	// TODO: rule for "sort" content like ASC or DESC
	if req.Sort == "" {
		req.Sort = "ASC"
	}

	switch {
	case req.Limit > 100:
		req.Limit = 100
	case req.Limit <= 0:
		req.Limit = 10
	}
	offset := (req.Page - 1) * req.Limit

	req.Offset = offset
	return
}
