package link

import (
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/gin-gonic/gin"
)

func (s service) HTTPAddLink(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeAddLink(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	link, err := s.AddLink(c.Request.Context(), dr)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := addLinkResponse{
		ID:       link.ID,
		URLShort: link.URLShort,
		URLFull:  link.URLFull,
		Title:    link.Title,
		Err:      err,
	}
	encodeResponse(c, sr)
}

func (s service) HTTPFindLinkByID(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeFindLinkByID(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	link, err := s.FindLinkByID(c.Request.Context(), dr)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := findLinkByIDResponse{
		Link: link,
		Err:  err,
	}
	encodeResponse(c, sr)
}

func (s service) HTTPFindLinks(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeFindLinks(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	links, err := s.FindLinks(c.Request.Context(), dr)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := findLinksResponse{
		Links: links,
		Err:   err,
	}
	encodeResponse(c, sr)
}

func (s service) HTTPUpdateLink(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeUpdateLink(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	link, err := s.UpdateLink(c.Request.Context(), dr)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := updateLinkResponse{
		Link: link,
		Err:  err,
	}
	encodeResponse(c, sr)
}

func (s service) HTTPDeleteLink(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeDeleteLink(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	err = s.DeleteLink(c.Request.Context(), dr)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := deleteLinkResponse{
		Err:  err,
	}
	encodeResponse(c, sr)
}

