package models

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

// Sub : define priority struct
type Sub struct {
	gorm.Model
	CategoryID  uint32   `json:"category_id"`
	Category    Category `gorm:"auto_preload" json:"category"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
}

// Prepare : prepare priority variables for save
func (sb *Sub) Prepare() {
	sb.Title = html.EscapeString(strings.TrimSpace(sb.Title))
	sb.Description = html.EscapeString(strings.TrimSpace(sb.Description))
}

// Validate : validate user input before save
func (sb *Sub) Validate() error {
	if sb.CategoryID == 0 {
		return errors.New("Required Category")
	}
	if sb.Title == "" {
		return errors.New("Required Title")
	}
	if sb.Description == "" {
		return errors.New("Required Description")
	}
	return nil
}

// SaveSub : create new priority
func (sb *Sub) SaveSub(db *gorm.DB) (*Sub, error) {
	var err error

	if categ := db.Where("id = ?", sb.CategoryID).First(&Category{}); categ.Error != nil {
		return &Sub{}, categ.Error
	}

	err = db.Debug().Model(&Sub{}).Create(&sb).Error
	if err != nil {
		return &Sub{}, err
	}
	if sb.ID != 0 {
		err = db.Debug().Model(&Category{}).Where("id = ?", sb.CategoryID).Take(&sb.Category).Error
		if err != nil {
			return &Sub{}, err
		}
	}
	return sb, nil
}

// FindAllSubs : get all subs
func (sb *Sub) FindAllSubs(db *gorm.DB) (*[]Sub, error) {
	var err error
	subs := []Sub{}
	err = db.Debug().Preload("Category").Preload("Category.Team").Limit(100).Find(&subs).Error
	if err != nil {
		return &[]Sub{}, err
	}
	return &subs, nil
}

// FindSubByID : find priority by priority ID
func (sb *Sub) FindSubByID(db *gorm.DB, cid uint64) (*Sub, error) {
	var err error
	// err = db.Debug().Model(&Sub{}).Where("id = ?", cid).Take(&sb).Error
	err = db.Debug().Where("id = ?", cid).Preload("Category").Preload("Category.Team").Find(&sb).Error
	if err != nil {
		return &Sub{}, err
	}
	return sb, nil
}

// UpdateASub : description
func (sb *Sub) UpdateASub(db *gorm.DB) (*Sub, error) {

	var err error

	err = db.Debug().Model(&Sub{}).Where("id = ?", sb.ID).Updates(Sub{
		CategoryID:  sb.CategoryID,
		Title:       sb.Title,
		Description: sb.Description,
	}).Error
	if err != nil {
		return &Sub{}, err
	}
	return sb, nil
}

// DeleteASub : delete a priority
func (sb *Sub) DeleteASub(db *gorm.DB, cid uint32) (int64, error) {

	db = db.Debug().Model(&Sub{}).Where("id = ?", cid).Take(&Sub{}).Delete(&Sub{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil

}
