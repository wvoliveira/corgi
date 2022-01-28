package link

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
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
	sort := c.Query("sort")
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
