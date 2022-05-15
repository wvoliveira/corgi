package token

import (
	"time"
)

type refreshResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    time.Time `json:"expires_in"`
	Err          error     `json:"err,omitempty"`
}

func (r refreshResponse) Error() error { return r.Err }
