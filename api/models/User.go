package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Tchayo/gotickets/api/utils/filter"
	"github.com/badoux/checkmail"
	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/jinzhu/gorm"
)

// User : description
type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	TeamID    uint      `json:"team_id"`
	Team      Team      `gorm:"auto_preload" json:"team"`
	Email     string    `gorm:"size:140;not null;unique" json:"email"`
	Password  string    `gorm:"size:140;not null" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Hash : description
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword : description
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// BeforeSave : description
func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Prepare : description
func (u *User) Prepare() {
	u.ID = 0
	u.Team = Team{}
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

// Validate : description
func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	case "reteam":
		if u.TeamID < 1 {
			return errors.New("Required Team")
		}
		return nil
	case "login":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	default:
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	}

}

// SaveUser : description
func (u *User) SaveUser(db *gorm.DB) (*User, error) {

	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	// hide user password
	u.Password = ""
	return u, nil
}

// FindAllUsers : description
func (u *User) FindAllUsers(db *gorm.DB, f *filter.Filter) (*pagination.Paginator, error) {

	var page, limit int
	users := []User{}

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
	}, &users)

	return res, nil

}

// FindUserByID : description
func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err

}

// UpdateAUser : description
func (u *User) UpdateAUser(db *gorm.DB, uid uint32) (*User, error) {

	// To hash pass
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":   u.Password,
			"email":      u.Email,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}

	// Display updated user
	err = db.Debug().Preload("Team").Where("id = ?", uid).Find(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

// UpdateUserTeam : description
func (u *User) UpdateUserTeam(db *gorm.DB, uid uint32) (*User, error) {

	// To hash pass
	var err error

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Updates(
		User{
			TeamID: u.TeamID,
		},
	)

	if db.Error != nil {
		return &User{}, db.Error
	}

	// Display updated user
	err = db.Debug().Preload("Team").Where("id = ?", uid).Find(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

// DeleteAUser : description
func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil

}
