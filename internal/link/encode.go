package link

import (
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/gin-gonic/gin"
)

type addLinkResponse struct {
	ID       string `json:"id"`
	URLShort string `json:"url_short"`
	URLFull  string `json:"url_full"`
	Title    string `json:"title"`
	Err      error  `json:"err,omitempty"`
}

func (r addLinkResponse) Error() error { return r.Err }

type findLinkByIDResponse struct {
	Link entity.Link `json:"data,omitempty"`
	Err  error       `json:"error,omitempty"`
}

func (r findLinkByIDResponse) Error() error { return r.Err }

type findLinksResponse struct {
	Links []entity.Link `json:"data,omitempty"`
	Err   error         `json:"error,omitempty"`
}

func (r findLinksResponse) Error() error { return r.Err }

type updateLinkResponse struct {
	Link entity.Link `json:"data,omitempty"`
	Err error `json:"err,omitempty"`
}

func (r updateLinkResponse) Error() error { return r.Err }

type deleteLinkResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteLinkResponse) Error() error { return r.Err }

func encodeResponse(c *gin.Context, response interface{}) {
	if err, ok := response.(e.Errorer); ok && err.Error() != nil {
		/*
			Not a Go kit transport error, but a business-logic error.
			Provide those as HTTP errors.
		*/
		e.EncodeError(c, err.Error())
	}
	c.JSON(200, response)
}
