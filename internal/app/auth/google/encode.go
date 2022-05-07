package google

import "time"

type loginResponse struct {
	RedirectURL string `json:"redirect_url"`
	Err         error  `json:"err,omitempty"`
}

func (r loginResponse) Error() error { return r.Err }

type callbackResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    time.Time `json:"expires_in"`
	Err          error     `json:"err,omitempty"`
}

func (r callbackResponse) Error() error { return r.Err }
