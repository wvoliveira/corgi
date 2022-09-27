package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Group struct {
	ID        string     `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`

	CreatedBy string `json:"created_by"`
	Users     []User `json:"users,omitempty" gorm:"many2many:user_groups;"`
}

func (l *Group) BeforeCreate(tx *gorm.DB) (err error) {
	l.ID = uuid.New().String()
	l.CreatedAt = time.Now()
	return
}

func (l *Group) BeforeUpdate(tx *gorm.DB) (err error) {
	t := time.Now()
	l.UpdatedAt = &t
	return
}
