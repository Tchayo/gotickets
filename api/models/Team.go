package models

import (
	"errors"
	"html"
	"strings"

	"github.com/Tchayo/gotickets/api/utils/filter"
	"github.com/badoux/checkmail"
	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/jinzhu/gorm"
)

// Team : define team struct
type Team struct {
	gorm.Model
	Name    string `json:"name"`
	Company string `json:"company"`
	Email   string `json:"email"`
}

// Prepare : prepare team variables for save
func (tm *Team) Prepare() {
	tm.Name = html.EscapeString(strings.TrimSpace(tm.Name))
	tm.Company = html.EscapeString(strings.TrimSpace(tm.Company))
	tm.Email = html.EscapeString(strings.TrimSpace(tm.Email))
}

// Validate : validate user input before save
func (tm *Team) Validate() error {
	if tm.Name == "" {
		return errors.New("Required Name")
	}
	if tm.Company == "" {
		return errors.New("Required Company")
	}
	if tm.Email == "" {
		return errors.New("Required Email")
	}
	if err := checkmail.ValidateFormat(tm.Email); err != nil {
		return errors.New("Invalid Email")
	}

	return nil
}

// SaveTeam : create new team
func (tm *Team) SaveTeam(db *gorm.DB) (*Team, error) {
	var err error
	err = db.Debug().Model(&Team{}).Create(&tm).Error
	if err != nil {
		return &Team{}, err
	}
	return tm, nil
}

// FindAllTeams : get all teams
func (tm *Team) FindAllTeams(db *gorm.DB, f *filter.Filter) (*pagination.Paginator, error) {

	var page, limit int
	teams := []Team{}

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
		DB:      db,
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"id desc"},
	}, &teams)

	return res, nil

}

// FindTeamByID : find team by team ID
func (tm *Team) FindTeamByID(db *gorm.DB, tmid uint64) (*Team, error) {
	var err error
	err = db.Debug().Model(&Team{}).Where("id = ?", tmid).Take(&tm).Error
	if err != nil {
		return &Team{}, err
	}
	return tm, nil
}

// UpdateATeam : description
func (tm *Team) UpdateATeam(db *gorm.DB) (*Team, error) {

	var err error

	err = db.Debug().Model(&Team{}).Where("id = ?", tm.ID).Updates(Team{
		Name:    tm.Name,
		Company: tm.Company,
		Email:   tm.Email,
	}).Error
	if err != nil {
		return &Team{}, err
	}
	return tm, nil
}

// DeleteATeam : delete a team
func (tm *Team) DeleteATeam(db *gorm.DB, tmid uint32) (int64, error) {

	db = db.Debug().Model(&Team{}).Where("id = ?", tmid).Take(&Team{}).Delete(&Team{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil

}
