package model

import (
	"encoding/json"
	"time"
)

type Session struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	TokenAccess  string    `json:"token_access"`
	TokenRefresh string    `json:"token_refresh"`
	ExpiresIn    time.Time `json:"expires_in"`

	User User `json:"user"`
}

func (s Session) Encode() (value []byte) {
	value, _ = json.Marshal(s)
	return
}
