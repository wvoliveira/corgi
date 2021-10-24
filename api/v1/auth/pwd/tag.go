package pwd

import (
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-pg/pg/orm"

	"github.com/elga-io/short/auth/jwt"
)

// Tag represents an authenticated application user
type Tag struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	AccountID int       `json:"-"`

	Name   string `json:"name"`
	URL    string `json:"url"`
	Active bool   `sql:",notnull" json:"active"`
}

// BeforeInsert hook executed before database insert operation.
func (a *Tag) BeforeInsert(db orm.DB) error {
	now := time.Now()
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
		a.UpdatedAt = now
	}
	return a.Validate()
}

// BeforeUpdate hook executed before database update operation.
func (a *Tag) BeforeUpdate(db orm.DB) error {
	a.UpdatedAt = time.Now()
	return a.Validate()
}

// BeforeDelete hook executed before database delete operation.
func (a *Tag) BeforeDelete(db orm.DB) error {
	return nil
}

// Validate validates Tag struct and returns validation errors.
func (a *Tag) Validate() error {
	a.Name = strings.TrimSpace(a.Name)

	return validation.ValidateStruct(a,
		validation.Field(&a.Name, validation.Required),
	)
}

// CanLogin returns true if user is allowed to login.
func (a *Tag) CanLogin() bool {
	return a.Active
}

// Claims returns the tag's claims to be signed
func (a *Tag) Claims() jwt.AppClaims {
	return jwt.AppClaims{
		ID: a.ID,
	}
}
