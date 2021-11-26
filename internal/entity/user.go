package entity

import "time"

// User represents a user info.
type User struct {
	ID        string    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	//LastLogin time.Time `json:"last_login"`

	Name string `json:"name"`
	//Email    string `json:"email" gorm:"index"`
	//Password string `json:"password"`

	Role string `json:"role"`
	Tags string `json:"tags"`

	Active string `json:"active"`

	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`

	Identities []Identity
	Tokens     []Token
	Links      []Link
}

// GetID returns the user ID.
func (u User) GetID() string {
	return u.ID
}

// GetEmail returns the e-mail.
//func (u User) GetEmail() string {
//	return u.Email
//}

// GetRole returns the role.
func (u User) GetRole() string {
	return u.Role
}
