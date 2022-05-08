package password

type registerResponse struct {
	Err error `json:"err,omitempty"`
}

func (r registerResponse) Error() error { return r.Err }
