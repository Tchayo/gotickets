package models

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

// Status : define priority struct
type Status struct {
	gorm.Model
	Title       string `json:"title"`
	Color       string `json:"color"`
	Description string `json:"description"`
	Resolution  bool   `gorm:"default:false" json:"resolution"`
	Lockable    bool   `gorm:"default:false" json:"lockable"`
	Assignable  bool   `gorm:"default:false" json:"assignable"`
}

// Prepare : prepare priority variables for save
func (st *Status) Prepare() {
	st.Title = html.EscapeString(strings.TrimSpace(st.Title))
	st.Color = html.EscapeString(strings.TrimSpace(st.Color))
	st.Description = html.EscapeString(strings.TrimSpace(st.Description))
}

// Validate : validate user input before save
func (st *Status) Validate() error {
	if st.Title == "" {
		return errors.New("Required Title")
	}
	if st.Color == "" {
		return errors.New("Required Color")
	}
	if st.Description == "" {
		return errors.New("Required Description")
	}
	return nil
}

// SaveStatus : create new priority
func (st *Status) SaveStatus(db *gorm.DB) (*Status, error) {
	var err error
	err = db.Debug().Model(&Status{}).Create(&st).Error
	if err != nil {
		return &Status{}, err
	}
	return st, nil
}

// FindAllStatuses : get all statuses
func (st *Status) FindAllStatuses(db *gorm.DB) (*[]Status, error) {
	var err error
	statuses := []Status{}
	err = db.Debug().Model(&Status{}).Limit(100).Find(&statuses).Error
	if err != nil {
		return &[]Status{}, err
	}
	return &statuses, nil
}

// FindStatusByID : find priority by priority ID
func (st *Status) FindStatusByID(db *gorm.DB, sid uint64) (*Status, error) {
	var err error
	err = db.Debug().Model(&Status{}).Where("id = ?", sid).Take(&st).Error
	if err != nil {
		return &Status{}, err
	}
	return st, nil
}

// UpdateAStatus : description
func (st *Status) UpdateAStatus(db *gorm.DB) (*Status, error) {

	var err error

	err = db.Debug().Model(&Status{}).Where("id = ?", st.ID).Updates(Status{
		Title:       st.Title,
		Color:       st.Color,
		Description: st.Description,
		Resolution:  st.Resolution,
		Lockable:    st.Lockable,
		Assignable:  st.Assignable,
	}).Error
	if err != nil {
		return &Status{}, err
	}
	return st, nil
}

// DeleteAStatus : delete a priority
func (st *Status) DeleteAStatus(db *gorm.DB, sid uint32) (int64, error) {

	db = db.Debug().Model(&Status{}).Where("id = ?", sid).Take(&Status{}).Delete(&Status{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil

}
