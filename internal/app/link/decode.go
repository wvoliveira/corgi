package link

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

type redirectRequest struct {
	WhoID   string
	Keyword string `uri:"keyword" binding:"required"`
	Domain  string
}

type addRequest struct {
	WhoID   string
	Domain  string `json:"domain"`
	Keyword string `json:"keyword"`
	URL     string `json:"url" binding:"required"`
	Title   string `json:"title"`
}

type findByIDRequest struct {
	WhoID  string
	LinkID string `uri:"id" binding:"required"`
}

type findAllRequest struct {
	WhoID        string
	Page         int
	Sort         string
	Offset       int
	Limit        int
	ShortenedURL string
	SearchText   string
}

type updateRequest struct {
	WhoID  string
	LinkID string `uri:"id" binding:"required"`
	Title  string `json:"title" binding:"required"`
}

type deleteRequest struct {
	WhoID   string
	LinkID  string `uri:"id" binding:"required"`
	Domain  string `json:"domain"`
	Keyword string `json:"keyword"`
}

type findFullURLRequest struct {
	WhoID   string
	Keyword string `uri:"keyword" binding:"required"`
	Domain  string
}

type clicksRequest struct {
	WhoID         string
	ShortURL      string
	TimestampFrom string
	TimestampTo   string
}

func decodeAdd(c *gin.Context) (req addRequest, err error) {
	v, ok := c.Get("user_id")
	if !ok {
		err = errors.New("impossible to know who you are")
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		return req, err
	}

	if req.Domain == "" {
		req.Domain = c.Request.Host
	}

	req.WhoID = v.(string)
	return req, nil
}

func decodeFindByID(c *gin.Context) (req findByIDRequest, err error) {
	v, ok := c.Get("user_id")
	if !ok {
		err = errors.New("impossible to know who you are")
		return
	}

	err = c.ShouldBindUri(&req)
	if err != nil {
		return req, errors.New("impossible to get link id from path")
	}

	req.WhoID = v.(string)
	return req, nil
}

func decodeFindAll(c *gin.Context) (req findAllRequest, err error) {
	v, ok := c.Get("user_id")
	if !ok {
		err = errors.New("impossible to know who you are")
		return
	}

	page := 1
	limit := 10

	queryPage := c.DefaultQuery("page", "1")
	queryLimit := c.DefaultQuery("limit", "10")

	if p, err := strconv.Atoi(queryPage); err == nil {
		page = p
	}

	if l, err := strconv.Atoi(queryLimit); err == nil {
		limit = l
		switch {
		case limit > 100:
			limit = 100
		case limit <= 0:
			limit = 10
		}
	}

	offset := (page - 1) * limit

	req.WhoID = v.(string)
	req.Page = page
	req.Sort = "id ASC"
	req.Limit = limit
	req.Offset = offset
	req.ShortenedURL = c.Query("u")
	req.SearchText = c.Query("q")
	return req, nil
}

func decodeUpdate(c *gin.Context) (req updateRequest, err error) {
	v, ok := c.Get("user_id")
	if !ok {
		err = errors.New("impossible to know who you are")
		return
	}

	if err = c.ShouldBindUri(&req); err != nil {
		return req, err
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		return req, err
	}

	req.WhoID = v.(string)
	return req, nil
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
	return req, nil
}

func decodeFindByKeyword(c *gin.Context) (req findFullURLRequest, err error) {
	v, ok := c.Get("user_id")
	if !ok {
		err = errors.New("impossible to know who you are")
		return
	}

	if err = c.ShouldBindUri(&req); err != nil {
		return req, err
	}

	req.WhoID = v.(string)
	req.Domain = c.Request.Host
	return req, nil
}

func decodeClicks(ctx *gin.Context) (req clicksRequest, err error) {
	shortURL := ctx.Query("u")
	timestampFrom := ctx.Query("tsf")
	timestampTo := ctx.Query("tst")

	if shortURL == "" {
		return req, errors.New("you need pass short URL with 'u' query param")
	}

	req.ShortURL = shortURL
	req.TimestampFrom = timestampFrom
	req.TimestampTo = timestampTo
	return
}
