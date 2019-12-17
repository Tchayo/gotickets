package models

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

// Priority : define priority struct
type Priority struct {
	gorm.Model
	TeamID      uint32 `json:"team_id"`
	Team        Team   `json:"team"`
	Title       string `json:"title"`
	Color       string `json:"color"`
	Description string `json:"description"`
	Fin         bool   `gorm:"default:false" json:"fin"`
}

// Prepare : prepare priority variables for save
func (pr *Priority) Prepare() {
	pr.Title = html.EscapeString(strings.TrimSpace(pr.Title))
	pr.Color = html.EscapeString(strings.TrimSpace(pr.Color))
	pr.Description = html.EscapeString(strings.TrimSpace(pr.Description))
}

// Validate : validate user input before save
func (pr *Priority) Validate() error {
	if pr.TeamID == 0 {
		return errors.New("Required Team")
	}
	if pr.Title == "" {
		return errors.New("Required Title")
	}
	if pr.Color == "" {
		return errors.New("Required Color")
	}
	if pr.Description == "" {
		return errors.New("Required Description")
	}
	return nil
}

// SavePriority : create new priority
func (pr *Priority) SavePriority(db *gorm.DB) (*Priority, error) {
	var err error
	err = db.Debug().Model(&Priority{}).Create(&pr).Error
	if err != nil {
		return &Priority{}, err
	}
	if pr.ID != 0 {
		err = db.Debug().Model(&Team{}).Where("id = ?", pr.TeamID).Take(&pr.Team).Error
		if err != nil {
			return &Priority{}, err
		}
	}
	return pr, nil
}

// FindAllPriorities : get all priorities
func (pr *Priority) FindAllPriorities(db *gorm.DB) (*[]Priority, error) {
	var err error
	priorities := []Priority{}
	err = db.Debug().Model(&Priority{}).Limit(100).Find(&priorities).Error
	if err != nil {
		return &[]Priority{}, err
	}
	if len(priorities) > 0 {
		for i := range priorities {
			err := db.Debug().Model(&User{}).Where("id = ?", priorities[i].TeamID).Take(&priorities[i].Team).Error
			if err != nil {
				return &[]Priority{}, err
			}
		}
	}
	return &priorities, nil
}

// FindPriorityByID : find priority by priority ID
func (pr *Priority) FindPriorityByID(db *gorm.DB, prid uint64) (*Priority, error) {
	var err error
	err = db.Debug().Model(&Priority{}).Where("id = ?", prid).Take(&pr).Error
	if err != nil {
		return &Priority{}, err
	}
	if pr.ID != 0 {
		err = db.Debug().Model(&Team{}).Where("id = ?", pr.TeamID).Take(&pr.Team).Error
		if err != nil {
			return &Priority{}, err
		}
	}
	return pr, nil
}

// UpdateAPriority : description
func (pr *Priority) UpdateAPriority(db *gorm.DB) (*Priority, error) {

	var err error

	err = db.Debug().Model(&Ticket{}).Where("id = ?", pr.ID).Updates(Priority{
		TeamID:      pr.TeamID,
		Title:       pr.Title,
		Color:       pr.Color,
		Description: pr.Description,
		Fin:         pr.Fin,
	}).Error
	if err != nil {
		return &Priority{}, err
	}
	return pr, nil
}

// DeleteAPriority : delete a priority
func (pr *Priority) DeleteAPriority(db *gorm.DB, prid uint32) (int64, error) {

	db = db.Debug().Model(&Priority{}).Where("id = ?", prid).Take(&Priority{}).Delete(&Priority{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil

}
