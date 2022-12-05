package link

import (
	"net/http"

	"github.com/gin-gonic/gin"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/middleware"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func (s service) NewHTTP(rg *gin.RouterGroup) {
	r := rg.Group("/links")
	r.Use(middleware.Checks())
	r.Use(middleware.Auth())

	r.OPTIONS("", nil)
	r.POST("", s.HTTPAdd)
	r.GET("/:id", s.HTTPFindByID)
	r.GET("/status/:id", s.HTTPFindByID)
	r.GET("", s.HTTPFindAll)
	r.PATCH("/:id", s.HTTPUpdate)
	r.DELETE("/:id", s.HTTPDelete)
}

func (s service) HTTPAdd(c *gin.Context) {

	d, err := decodeAdd(c)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	link, err := s.Add(c, model.Link{
		Domain:  d.Domain,
		Keyword: d.Keyword,
		URL:     d.URL,
		Title:   d.Title,
		UserID:  d.UserID,
	})

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	response.Default(c, link, "", http.StatusCreated)
}

func (s service) HTTPFindByID(c *gin.Context) {

	d, err := decodeFindByID(c)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	link, err := s.FindByID(c, d.ID, d.UserID)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	response.Default(c, link, "", http.StatusOK)
}

func (s service) HTTPFindAll(c *gin.Context) {

	dr, err := decodeFindAll(c)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	total, pages, links, err := s.FindAll(c, dr.Offset, dr.Limit, dr.Sort, dr.UserID)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	sr := findAllResponse{
		Links: links,
		Limit: dr.Limit,
		Page:  dr.Page,
		Sort:  dr.Sort,
		Total: total,
		Pages: pages,
		Err:   err,
	}

	response.Default(c, sr, "", http.StatusOK)
}

func (s service) HTTPUpdate(c *gin.Context) {

	dr, userID, err := decodeUpdate(c)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	link, err := s.Update(c, model.Link{
		ID:     dr.ID,
		Title:  dr.Title,
		UserID: userID,
	})

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	response.Default(c, link, "", http.StatusOK)
}

func (s service) HTTPDelete(c *gin.Context) {

	dr, err := decodeDelete(c)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	err = s.Delete(c, dr.ID, dr.UserID)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	response.Default(c, nil, "", http.StatusOK)
}
