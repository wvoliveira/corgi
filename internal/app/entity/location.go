package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LocationIPv4 get info location from database by IPv4 or IPv6 address.
type LocationIPv4 struct {
	ID         string  `json:"id" gorm:"primaryKey;autoIncrement:false"`
	RangeStart uint32  `json:"range_start" gorm:"index;unique"`
	RangeEnd   uint32  `json:"range_end" gorm:"index;unique"`
	Country    string  `json:"country"`
	State      string  `json:"state"`
	City       string  `json:"city"`
	Latitude   float64 `json:"latitude" gorm:"type:decimal(10,3);"`
	Longitude  float64 `json:"longitude" gorm:"type:decimal(10,3);"`
}

// LocationIPv6 get info location from database by IPv4 or IPv6 address. // gorm:"index"
type LocationIPv6 struct {
	ID         string  `json:"id" gorm:"primaryKey;autoIncrement:false"`
	RangeStart string  `json:"range_start" gorm:"index;type:BIGINT;unique"`
	RangeEnd   string  `json:"range_end" gorm:"index;type:BIGINT;unique"`
	Country    string  `json:"country"`
	State      string  `json:"state"`
	City       string  `json:"city"`
	Latitude   float64 `json:"latitude" gorm:"type:decimal(10,3);"`
	Longitude  float64 `json:"longitude" gorm:"type:decimal(10,3);"`
}

func (l *LocationIPv4) BeforeCreate(tx *gorm.DB) (err error) {
	l.ID = uuid.New().String()
	return
}
