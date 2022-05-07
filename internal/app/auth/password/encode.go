package password

import (
	"net/http"
	"time"

	e "github.com/elga-io/corgi/internal/pkg/errors"
)

type loginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    time.Time `json:"expires_in"`
	Err          error     `json:"err,omitempty"`
}

func (r loginResponse) Error() error { return r.Err }

type registerResponse struct {
	Err error `json:"err,omitempty"`
}

func (r registerResponse) Error() error { return r.Err }

func encodeResponse(w http.ResponseWriter, response interface{}) {
	if err, ok := response.(e.Errorer); ok && err.Error() != nil {
		e.EncodeError(w, err.Error())
		return
	}
	c.JSON(200, response)
}
