package models

import (
	"errors"
	"html"
	"strings"

	"github.com/Tchayo/gotickets/api/utils/filter"
	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/jinzhu/gorm"
)

// Priority : define priority struct
type Priority struct {
	gorm.Model
	TeamID         uint   `json:"team_id"`
	Team           Team   `gorm:"auto_preload" json:"team"`
	Title          string `json:"title"`
	TimeoutHours   int    `gorm:"default:0" json:"timeout_hours"`
	TimeoutMinutes int    `gorm:"default:0" json:"timeout_minutes"`
	Color          string `json:"color"`
	Description    string `json:"description"`
	Fin            bool   `gorm:"default:false" json:"fin"`
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

	if prteam := db.Where("id = ?", pr.TeamID).First(&Team{}); prteam.Error != nil {
		return &Priority{}, prteam.Error
	}

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
func (pr *Priority) FindAllPriorities(db *gorm.DB, f *filter.Filter) (*pagination.Paginator, error) {

	var page, limit int
	priorities := []Priority{}

	if f.Page < 1 {
		page = 1
	} else {
		page = f.Page
	}
	if f.Limit == 0 {
		limit = 10
	} else {
		limit = f.Limit
	}

	res := pagination.Paging(&pagination.Param{
		DB:      db.Preload("Team"),
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"id desc"},
	}, &priorities)

	return res, nil

}

// FindPriorityByID : find priority by priority ID
func (pr *Priority) FindPriorityByID(db *gorm.DB, prid uint64) (*Priority, error) {
	var err error
	err = db.Debug().Preload("Team").Where("id = ?", prid).Find(&pr).Error
	if err != nil {
		return &Priority{}, err
	}
	return pr, nil
}

// UpdateAPriority : description
func (pr *Priority) UpdateAPriority(db *gorm.DB) (*Priority, error) {

	var err error

	err = db.Debug().Model(&Priority{}).Where("id = ?", pr.ID).Updates(Priority{
		TeamID:         pr.TeamID,
		Title:          pr.Title,
		TimeoutHours:   pr.TimeoutHours,
		TimeoutMinutes: pr.TimeoutMinutes,
		Color:          pr.Color,
		Description:    pr.Description,
		Fin:            pr.Fin,
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
