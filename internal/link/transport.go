package link

import (
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func (s service) Routers(e *gin.Engine) {
	r := e.Group("/api/v1/links",
		middlewares.Checks(s.logger),
		middlewares.Auth(s.logger, s.secret))

	r.OPTIONS("/", nil)
	r.POST("/", s.HTTPAdd)
	r.GET("/:id", s.HTTPFindByID)
	r.GET("/", s.HTTPFindAll)
	r.PATCH("/:id", s.HTTPUpdate)
	r.DELETE("/:id", s.HTTPDelete)
}

func (s service) HTTPAdd(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeAdd(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	link, err := s.Add(c.Request.Context(), entity.Link{Domain: dr.Domain, Keyword: dr.Keyword, URL: dr.URL, Title: dr.Title, UserID: dr.UserID})
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := addResponse{
		ID:      link.ID,
		Domain:  link.Domain,
		Keyword: link.Keyword,
		URL:     link.URL,
		Title:   link.Title,
		Err:     err,
	}
	encodeResponse(c, sr)
}

func (s service) HTTPFindByID(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeFindByID(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	link, err := s.FindByID(c.Request.Context(), dr.ID, dr.UserID)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := findByIDResponse{
		Link: link,
		Err:  err,
	}
	encodeResponse(c, sr)
}

func (s service) HTTPFindAll(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeFindAll(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	total, pages, links, err := s.FindAll(c.Request.Context(), dr.Offset, dr.Limit, dr.Sort, dr.UserID)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := findAllResponse{
		Links: links,
		Limit: dr.Limit,
		Page:  dr.Page,
		Sort:  dr.Sort,
		Total: total,
		Pages: pages,
		Err:   err,
	}
	encodeResponse(c, sr)
}

func (s service) HTTPUpdate(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeUpdate(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	link, err := s.Update(
		c.Request.Context(),
		entity.Link{ID: dr.ID, Domain: dr.Domain, Keyword: dr.Keyword, URL: dr.URL, Title: dr.Title, Active: dr.Active, UserID: dr.UserID},
	)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := updateResponse{
		Link: link,
		Err:  err,
	}
	encodeResponse(c, sr)
}

func (s service) HTTPDelete(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeDelete(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	err = s.Delete(c.Request.Context(), dr.ID, dr.UserID)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := deleteResponse{
		Err: err,
	}
	encodeResponse(c, sr)
}
