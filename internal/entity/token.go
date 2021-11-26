package entity

import "time"

// Token represent a JWT struct.
type Token struct {
	ID        string    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastUse   time.Time `json:"last_use"`

	AccessToken  string    `json:"-" gorm:"-"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    time.Time `json:"expires_in"`

	UserID       string    `json:"user_id" gorm:"index"`
}

// GetID returns the Token ID.
func (t Token) GetID() string {
	return t.ID
}
