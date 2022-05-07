package link

import (
	"github.com/elga-io/corgi/internal/app/entity"
	e "github.com/elga-io/corgi/internal/pkg/errors"
	"github.com/elga-io/corgi/internal/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func (s service) NewHTTP(e *gin.Engine) {
	r := e.Group("/api/v1/links",
		middlewares.Checks(),
		middlewares.Auth(s.secret))

	r.OPTIONS("", nil)
	r.POST("", s.HTTPAdd)
	r.GET(":id", s.HTTPFindByID)
	r.GET("/status/:id", s.HTTPFindByID)
	r.GET("", s.HTTPFindAll)
	r.PATCH(":id", s.HTTPUpdate)
	r.DELETE(":id", s.HTTPDelete)
}

func (s service) HTTPAdd(c *gin.Context) {
	dr, err := decodeAdd(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	link, err := s.Add(c, entity.Link{Domain: dr.Domain, Keyword: dr.Keyword, URL: dr.URL, Title: dr.Title, UserID: dr.UserID})
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	encodeResponse(c, link)
}

func (s service) HTTPFindByID(c *gin.Context) {
	dr, err := decodeFindByID(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	link, err := s.FindByID(c, dr.ID, dr.UserID)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	encodeResponse(c, link)
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

	encodeResponse(c, findAllResponse{
		Links: links,
		Limit: dr.Limit,
		Page:  dr.Page,
		Sort:  dr.Sort,
		Total: total,
		Pages: pages,
		Err:   err,
	})
}

func (s service) HTTPUpdate(c *gin.Context) {
	dr, err := decodeUpdate(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	link, err := s.Update(c, entity.Link{
		ID:      dr.ID,
		Domain:  dr.Domain,
		Keyword: dr.Keyword,
		URL:     dr.URL,
		Title:   dr.Title,
		Active:  dr.Active,
		UserID:  dr.UserID,
	})
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	encodeResponse(c, link)
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

	encodeResponse(c, nil)
}
