package link

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/pkg/model"
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
	Page         int    `json:"page"`
	Sort         string `json:"sort"`
	Offset       int    `json:"offset"`
	Limit        int    `json:"limit"`
	UserID       string `json:"user_id"`
	ShortenedURL string
}

type updateRequest struct {
	ID        string    `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
	Title     string    `json:"title"`
}

type deleteRequest struct {
	ID      string `json:"id"`
	Domain  string `json:"domain"`
	Keyword string `json:"keyword"`
	UserID  string `json:"user_id"`
}

func decodeAdd(c *gin.Context) (r addRequest, err error) {

	session := sessions.Default(c)
	v := session.Get("user")

	if v == nil {
		return r, errors.New("impossible to get user from session")
	}

	if err = c.ShouldBindJSON(&r); err != nil {
		return r, err
	}

	// TODO: insert this logic in config env/file.
	// with default domain if not send in payload.
	if r.Domain == "" {
		r.Domain = c.Request.Host
	}

	r.UserID = v.(model.User).ID
	return r, nil
}

func decodeFindByID(c *gin.Context) (r findByIDRequest, err error) {

	session := sessions.Default(c)
	v := session.Get("user")

	if v == nil {
		return r, errors.New("impossible to get user from session")
	}

	linkID := c.Param("id")
	if linkID == "" {
		return r, errors.New("impossible to get link id from path")
	}

	r.ID = linkID
	r.UserID = v.(model.User).ID
	return r, nil
}

func decodeFindAll(c *gin.Context) (r findAllRequest, err error) {

	session := sessions.Default(c)
	v := session.Get("user")

	if v == nil {
		return r, errors.New("impossible to get user from session")
	}

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	sort := c.Query("sort")
	shortenedURL := c.Query("u")

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

	r.Page = page
	r.Sort = sort
	r.Limit = limit
	r.Offset = offset
	r.UserID = v.(model.User).ID
	r.ShortenedURL = shortenedURL

	return r, nil
}

func decodeUpdate(c *gin.Context) (r updateRequest, userID string, err error) {

	session := sessions.Default(c)
	v := session.Get("user")

	if v == nil {
		return r, userID, errors.New("impossible to get user from session")
	}

	linkID := c.Param("id")
	if linkID == "" {
		return r, userID, errors.New("impossible to get link id from path")
	}

	if err = c.ShouldBindJSON(&r); err != nil {
		return r, userID, err
	}

	r.ID = linkID
	userID = v.(model.User).ID

	return r, userID, nil
}

func decodeDelete(c *gin.Context) (r deleteRequest, err error) {

	session := sessions.Default(c)
	v := session.Get("user")

	if v == nil {
		return r, errors.New("impossible to get user from session")
	}

	linkID := c.Param("id")
	if linkID == "" {
		return r, errors.New("impossible to get link id from path")
	}

	r.ID = linkID
	r.UserID = v.(model.User).ID

	return r, nil
}
