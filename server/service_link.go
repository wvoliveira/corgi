package server

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AddLink create a new short Link.
func (s Service) AddLink(auth Account, payload Link) (link Link, err error) {
	err = s.db.Model(&Link{}).Where("keyword = ?", payload.Keyword).Take(&link).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		payload.ID = uuid.New().String()
		payload.CreatedAt = time.Now()
		payload.AccountID = auth.ID
		payload.Active = "true"

		link = payload
		err = s.db.Model(&Link{}).Create(&link).Error
		return
	} else if err == nil {
		return link, ErrAlreadyExists
	}
	return
}

// FindLinkByID search a specific Link by ID.
func (s Service) FindLinkByID(auth Account, id string) (link Link, err error) {
	err = s.db.Model(&Link{}).Where("account_id = ?", auth.ID).Where("id = ?", id).First(&link).Error
	return
}

// FindLinks get a Link list from database.
func (s Service) FindLinks(auth Account, offset, limit int) (links []Link, err error) {
	err = s.db.Model(&Link{}).Where("account_id = ?", auth.ID).Limit(limit).Offset(offset).Find(links).Error
	return
}

// UpdateLink update specific Link by ID.
func (s Service) UpdateLink(auth Account, id string, payload Link) (err error) {
	err = s.db.Model(&Link{}).Where("account_id = ?", auth.ID).Where("id = ?", id).Updates(&payload).Error
	return
}

// DeleteLink delete Link by ID.
func (s Service) DeleteLink(auth Account, id string) (err error) {
	var link Link
	err = s.db.Model(&Link{}).Where("account_id = ?", auth.ID).Where("id = ?", id).Delete(&link).Error
	return
}
