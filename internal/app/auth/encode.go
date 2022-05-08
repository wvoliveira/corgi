package auth

type logoutResponse struct {
	Err error `json:"err,omitempty"`
}

func (r logoutResponse) Error() error { return r.Err }
