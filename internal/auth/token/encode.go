package token

import (
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/gin-gonic/gin"
	"time"
)

type refreshResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    time.Time `json:"expires_in"`
	Err          error     `json:"err,omitempty"`
}

func (r refreshResponse) Error() error { return r.Err }

func encodeResponse(c *gin.Context, response interface{}) {
	if err, ok := response.(e.Errorer); ok && err.Error() != nil {
		e.EncodeError(c, err.Error())
	}
	c.JSON(200, response)
}
