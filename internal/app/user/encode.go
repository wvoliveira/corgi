package user

import (
	"time"
)

type identity struct {
	Provider string `json:"provider"`
	UID      string `json:"uid"`
}

type userResponse struct {
	Name       string     `json:"name"`
	Role       string     `json:"role"`
	Identities []identity `json:"identities"`
}

type updateResponse struct {
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Err       error     `json:"err,omitempty"`
}

func (r updateResponse) Error() error { return r.Err }
