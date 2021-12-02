package link

import (
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/gin-gonic/gin"
)

type addResponse struct {
	ID       string `json:"id"`
	URLShort string `json:"url_short"`
	URLFull  string `json:"url_full"`
	Title    string `json:"title"`
	Err      error  `json:"err,omitempty"`
}

func (r addResponse) Error() error { return r.Err }

type findByIDResponse struct {
	Link entity.Link `json:"data,omitempty"`
	Err  error       `json:"error,omitempty"`
}

func (r findByIDResponse) Error() error { return r.Err }

type findAllResponse struct {
	Links []entity.Link `json:"data,omitempty"`
	Err   error         `json:"error,omitempty"`
}

func (r findAllResponse) Error() error { return r.Err }

type updateResponse struct {
	Link entity.Link `json:"data,omitempty"`
	Err  error       `json:"err,omitempty"`
}

func (r updateResponse) Error() error { return r.Err }

type deleteResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteResponse) Error() error { return r.Err }

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
