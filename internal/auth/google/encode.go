package google

import (
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/gin-gonic/gin"
)

type loginResponse struct {
	RedirectURL string `json:"redirect_url"`
	Err         error  `json:"err,omitempty"`
}

func (r loginResponse) Error() error { return r.Err }

type callbackResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Err          error  `json:"err,omitempty"`
}

func (r callbackResponse) Error() error { return r.Err }

func encodeResponse(c *gin.Context, response interface{}) {
	if err, ok := response.(e.Errorer); ok && err.Error() != nil {
		e.EncodeError(c, err.Error())
		return
	}
	c.JSON(200, response)
}
