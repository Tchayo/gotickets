package models

import (
	"errors"
	"html"
	"strings"

	"github.com/Tchayo/gotickets/api/utils/filter"
	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/jinzhu/gorm"
)

// Category : define priority struct
type Category struct {
	gorm.Model
	TeamID      uint   `json:"team_id"`
	Team        Team   `gorm:"auto_preload" json:"team"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Prepare : prepare priority variables for save
func (ct *Category) Prepare() {
	ct.Title = html.EscapeString(strings.TrimSpace(ct.Title))
	ct.Description = html.EscapeString(strings.TrimSpace(ct.Description))
}

// Validate : validate user input before save
func (ct *Category) Validate() error {
	if ct.TeamID == 0 {
		return errors.New("Required Team")
	}
	if ct.Title == "" {
		return errors.New("Required Title")
	}
	if ct.Description == "" {
		return errors.New("Required Description")
	}
	return nil
}

// SaveCategory : create new priority
func (ct *Category) SaveCategory(db *gorm.DB) (*Category, error) {
	var err error

	if cteam := db.Where("id = ?", ct.TeamID).First(&Team{}); cteam.Error != nil {
		return &Category{}, cteam.Error
	}

	err = db.Debug().Model(&Category{}).Create(&ct).Error
	if err != nil {
		return &Category{}, err
	}
	if ct.ID != 0 {
		err = db.Debug().Model(&Team{}).Where("id = ?", ct.TeamID).Take(&ct.Team).Error
		if err != nil {
			return &Category{}, err
		}
	}
	return ct, nil
}

// FindAllCategories : get all categories
func (ct *Category) FindAllCategories(db *gorm.DB, f *filter.Filter) (*pagination.Paginator, error) {

	var page, limit int
	categories := []Category{}

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
	}, &categories)

	return res, nil

}

// FindCategoryByID : find priority by priority ID
func (ct *Category) FindCategoryByID(db *gorm.DB, cid uint64) (*Category, error) {
	var err error
	err = db.Debug().Where("id = ?", cid).Preload("Team").Find(&ct).Error
	if err != nil {
		return &Category{}, err
	}
	return ct, nil
}

// UpdateACategory : description
func (ct *Category) UpdateACategory(db *gorm.DB) (*Category, error) {

	var err error

	err = db.Debug().Model(&Category{}).Where("id = ?", ct.ID).Updates(Category{
		TeamID:      ct.TeamID,
		Title:       ct.Title,
		Description: ct.Description,
	}).Error
	if err != nil {
		return &Category{}, err
	}
	return ct, nil
}

// DeleteACategory : delete a priority
func (ct *Category) DeleteACategory(db *gorm.DB, cid uint32) (int64, error) {

	db = db.Debug().Model(&Category{}).Where("id = ?", cid).Take(&Category{}).Delete(&Category{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil

}
