package link

import (
	"net/http"

	"github.com/gin-gonic/gin"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/middleware"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

// NewHTTP create a new http endpoint for link service.
// This is a principal service, and it has a root router to redirect
// links by domain and keyword combination.
func (s service) NewHTTP(root *gin.Engine, rg *gin.RouterGroup) {
	root.GET("/:keyword", s.HTTPRedirect)

	r := rg.Group("/links")
	r.Use(middleware.Checks())

	r.POST("", s.HTTPAdd)
	r.GET("", s.HTTPFindAll)
	r.GET("/:id", s.HTTPFindByID)
	r.PATCH("/:id", s.HTTPUpdate)
	r.DELETE("/:id", s.HTTPDelete)
	r.GET("/keyword/:keyword", s.HTTPFindFullURL)
	//r.GET("/clicks", s.HTTPClicks)
	r.GET("/clicks", s.HTTPClicks)
}

func (s service) HTTPRedirect(ctx *gin.Context) {
	d, err := decodeFindByKeyword(ctx)
	if err != nil {
		e.EncodeError(ctx, err)
		return
	}

	link, err := s.FindRedirectURL(ctx, d.Domain, d.Keyword)
	if err != nil {
		e.EncodeError(ctx, err)
		return
	}

	url := encodeRedirect(link)
	ctx.Redirect(301, url.URL)
}

func (s service) HTTPAdd(ctx *gin.Context) {
	payload, err := decodeAdd(ctx)
	if err != nil {
		e.EncodeError(ctx, err)
		return
	}

	link, err := s.Add(ctx, payload)
	if err != nil {
		e.EncodeError(ctx, err)
		return
	}

	response.Default(ctx, link, "", http.StatusCreated)
}

func (s service) HTTPFindByID(c *gin.Context) {
	payload, err := decodeFindByID(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	link, err := s.FindByID(c, payload)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	response.Default(c, link, "", http.StatusOK)
}

func (s service) HTTPFindAll(ctx *gin.Context) {
	payload, err := decodeFindAll(ctx)
	if err != nil {
		e.EncodeError(ctx, err)
		return
	}

	total, pages, links, err := s.FindAll(ctx, payload)
	if err != nil {
		e.EncodeError(ctx, err)
		return
	}

	sr := findAllResponse{
		Links: links,
		Limit: payload.Limit,
		Page:  payload.Page,
		Sort:  payload.Sort,
		Total: total,
		Pages: pages,
		Err:   err,
	}

	response.Default(ctx, sr, "", http.StatusOK)
}

func (s service) HTTPUpdate(ctx *gin.Context) {
	payload, err := decodeUpdate(ctx)
	if err != nil {
		e.EncodeError(ctx, err)
		return
	}

	err = s.Update(ctx, payload)
	if err != nil {
		e.EncodeError(ctx, err)
		return
	}

	response.Default(ctx, nil, "", http.StatusOK)
}

func (s service) HTTPDelete(ctx *gin.Context) {
	payload, err := decodeDelete(ctx)
	if err != nil {
		e.EncodeError(ctx, err)
		return
	}

	err = s.Delete(ctx, payload)
	if err != nil {
		e.EncodeError(ctx, err)
		return
	}

	response.Default(ctx, nil, "", http.StatusOK)
}

func (s service) HTTPFindFullURL(c *gin.Context) {
	d, err := decodeFindByKeyword(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	link, err := s.FindFullURL(c, d.Domain, d.Keyword)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	url := encodeFindByKeyword(link)
	response.Default(c, url, "", http.StatusOK)
}

func (s service) HTTPClicks(c *gin.Context) {
	payload, err := decodeClicks(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	linkClicks, err := s.Clicks(c, payload)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	//url := encodeFindByKeyword(linkClicks)
	response.Default(c, linkClicks, "", http.StatusOK)
}
