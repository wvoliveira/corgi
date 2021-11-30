package link

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
)

type addLinkRequest struct {
	URLShort string `json:"url_short"`
	URLFull  string `json:"url_full"`
	Title    string `json:"title"`
	UserID   string `json:"user_id"`
}

type findLinkByIDRequest struct {
	ID     string
	UserID string `json:"user_id"`
}

type findLinksRequest struct {
	Offset int
	Limit  int
	UserID string `json:"user_id"`
}

type updateLinkRequest struct {
	ID       string
	URLShort string `json:"url_short" gorm:"index"`
	URLFull  string `json:"url_full"`
	Title    string `json:"title"`
	Active   string `json:"active"`
	UserID   string `json:"user_id"`
}

type deleteLinkRequest struct {
	ID     string
	UserID string `json:"user_id"`
}

func decodeAddLink(c *gin.Context) (req addLinkRequest, err error) {
	userID, _ := c.Get("user_id")
	if err = json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return req, err
	}
	req.UserID = userID.(string)
	return req, nil
}

func decodeFindLinkByID(c *gin.Context) (req findLinkByIDRequest, err error) {
	userID, _ := c.Get("user_id")
	linkID := c.Param("id")

	req.ID = linkID
	req.UserID = userID.(string)
	return req, nil
}

func decodeFindLinks(c *gin.Context) (req findLinksRequest, err error) {
	userID, _ := c.Get("user_id")

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if page == 0 {
		page = 1
	}

	switch {
	case limit > 100:
		limit = 100
	case limit <= 0:
		limit = 10
	}
	offset := (page - 1) * limit

	req.Limit = limit
	req.Offset = offset
	req.UserID = userID.(string)
	return req, nil
}

func decodeUpdateLink(c *gin.Context) (req updateLinkRequest, err error) {
	userID, _ := c.Get("user_id")
	linkID := c.Param("id")

	if err = json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return req, err
	}
	req.ID = linkID
	req.UserID = userID.(string)
	return req, nil
}

func decodeDeleteLink(c *gin.Context) (req deleteLinkRequest, err error) {
	userID, _ := c.Get("user_id")
	linkID := c.Param("id")
	req.ID = linkID
	req.UserID = userID.(string)
	return req, nil
}
