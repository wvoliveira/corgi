package password

import (
	"encoding/json"
	"net/http"
	e "github.com/elga-io/corgi/pkg/errors"
)

type authLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Err          error  `json:"err,omitempty"`
}

func (r authLoginResponse) error() error { return r.Err }

type authRegisterResponse struct {
	Err error `json:"err,omitempty"`
}

func (r authRegisterResponse) error() error { return r.Err }

func encodeResponse(w http.ResponseWriter, response interface{}) error {
	if err, ok := response.(e.Errorer); ok && err.Error() != nil {
		/*
			Not a Go kit transport error, but a business-logic error.
			Provide those as HTTP errors.
		*/
		e.EncodeError(err.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
