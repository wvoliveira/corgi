package log

import (
	"time"
)

type identity struct {
	Provider string `json:"provider,omitempty"`
	UID      string `json:"uid,omitempty"`
}

type userResponse struct {
	Name       string     `json:"name"`
	Role       string     `json:"role,omitempty"`
	Identities []identity `json:"identities,omitempty"`
}

type updateResponse struct {
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Err       error     `json:"err,omitempty"`
}

func (r updateResponse) Error() error { return r.Err }
