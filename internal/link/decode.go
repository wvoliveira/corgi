package link

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type addRequest struct {
	URLShort string `json:"url_short"`
	URLFull  string `json:"url_full"`
	Title    string `json:"title"`
	UserID   string `json:"user_id"`
}

type findByIDRequest struct {
	ID     string
	UserID string `json:"user_id"`
}

type findAllRequest struct {
	Offset int
	Limit  int
	UserID string `json:"user_id"`
}

type updateRequest struct {
	ID        string    `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	URLShort  string    `json:"url_short"`
	URLFull   string    `json:"url_full"`
	Title     string    `json:"title"`
	Active    string    `json:"active"`
	UserID    string    `json:"-"`
}

type deleteRequest struct {
	ID       string
	URLShort string
	UserID   string `json:"user_id"`
}

func decodeAdd(c *gin.Context) (req addRequest, err error) {
	userID, _ := c.Get("user_id")
	if err = json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return req, err
	}
	req.UserID = userID.(string)
	return req, nil
}

func decodeFindByID(c *gin.Context) (req findByIDRequest, err error) {
	userID, _ := c.Get("user_id")
	linkID := c.Param("id")

	req.ID = linkID
	req.UserID = userID.(string)
	return req, nil
}

func decodeFindAll(c *gin.Context) (req findAllRequest, err error) {
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

func decodeUpdate(c *gin.Context) (req updateRequest, err error) {
	userID, _ := c.Get("user_id")
	linkID := c.Param("id")

	if err = json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return req, err
	}
	req.ID = linkID
	req.UserID = userID.(string)
	return req, nil
}

func decodeDelete(c *gin.Context) (req deleteRequest, err error) {
	userID, _ := c.Get("user_id")
	linkID := c.Param("id")
	req.ID = linkID
	req.UserID = userID.(string)
	return req, nil
}
