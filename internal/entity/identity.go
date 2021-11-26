package entity

import "time"

type Identity struct {
	ID        string    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastLogin time.Time `json:"last_login"`

	Provider    string    `json:"provider"` // phone, email, wechat, github...
	UID         string    `json:"uid"`      // e-mail, google id, facebook id, etc
	Password    string    `json:"password"`
	UserID      string    `json:"user_id"`
	ConfirmedAt time.Time `json:"confirmed_at"`
}

// GetID returns the user ID.
func (i Identity) GetID() string {
	return i.ID
}

// GetUID returns e-mail or google id or facebook id, etc.
func (i Identity) GetUID() string {
	return i.UID
}
