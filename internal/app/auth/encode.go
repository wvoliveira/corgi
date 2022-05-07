package auth

import (
	e "github.com/elga-io/corgi/internal/pkg/errors"
	"github.com/gin-gonic/gin"
)

type logoutResponse struct {
	Err error `json:"err,omitempty"`
}

func (r logoutResponse) Error() error { return r.Err }

func encodeResponse(c *gin.Context, response interface{}) {
	if err, ok := response.(e.Errorer); ok && err.Error() != nil {
		e.EncodeError(c, err.Error())
	}
	c.JSON(200, response)
}
