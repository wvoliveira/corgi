package password

import (
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/gin-gonic/gin"
)

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Err          error  `json:"err,omitempty"`
}

func (r loginResponse) Error() error { return r.Err }

type registerResponse struct {
	Err error `json:"err,omitempty"`
}

func (r registerResponse) Error() error { return r.Err }

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
