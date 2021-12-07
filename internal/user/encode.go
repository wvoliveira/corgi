package user

import (
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/gin-gonic/gin"
	"time"
)

type userResponse struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type findResponse struct {
	userResponse `json:"data,omitempty"`
	Err          error `json:"error,omitempty"`
}

func (r findResponse) Error() error { return r.Err }

type updateResponse struct {
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Err       error     `json:"err,omitempty"`
}

func (r updateResponse) Error() error { return r.Err }

func encodeResponse(c *gin.Context, response interface{}) {
	if err, ok := response.(e.Errorer); ok && err.Error() != nil {
		e.EncodeError(c, err.Error())
	}
	c.JSON(200, response)
}
