package link

import (
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/gin-gonic/gin"
)

func (s service) Routers(r *gin.RouterGroup) {
	r.POST("/links", s.HTTPAdd)
	r.GET("/links/:id", s.HTTPFindByID)
	r.GET("/links", s.HTTPFindAll)
	r.PATCH("/links/:id", s.HTTPUpdate)
	r.DELETE("/links/:id", s.HTTPDelete)
}

func (s service) HTTPAdd(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeAdd(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	link, err := s.Add(c.Request.Context(), dr)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := addResponse{
		ID:       link.ID,
		URLShort: link.URLShort,
		URLFull:  link.URLFull,
		Title:    link.Title,
		Err:      err,
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
	link, err := s.FindByID(c.Request.Context(), dr)
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
	links, err := s.FindAll(c.Request.Context(), dr)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := findAllResponse{
		Links: links,
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
	link, err := s.Update(c.Request.Context(), dr)
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
	err = s.Delete(c.Request.Context(), dr)
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
