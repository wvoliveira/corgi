package link

import (
	"github.com/elga-io/corgi/internal/app/entity"
	e "github.com/elga-io/corgi/internal/pkg/errors"
	"github.com/gin-gonic/gin"
)

type addResponse struct {
	ID      string `json:"id"`
	Domain  string `json:"domain"`
	Keyword string `json:"keyword"`
	URL     string `json:"url"`
	Title   string `json:"title"`
	Err     error  `json:"err,omitempty"`
}

func (r addResponse) Error() error { return r.Err }

type findByIDResponse struct {
	Link entity.Link `json:"data,omitempty"`
	Err  error       `json:"error,omitempty"`
}

func (r findByIDResponse) Error() error { return r.Err }

type findAllResponse struct {
	Links []entity.Link `json:"data"`
	Limit int           `json:"limit"`
	Page  int           `json:"page"`
	Sort  string        `json:"sort"`
	Total int64         `json:"total"`
	Pages int           `json:"pages"`
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
		e.EncodeError(c, err.Error())
	}
	c.JSON(200, response)
}
