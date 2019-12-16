package models

import (
	"errors"
	"fmt"
	"html"
	"os"
	"strings"

	"github.com/dongri/phonenumber"
	"github.com/jinzhu/gorm"
	"github.com/speps/go-hashids"
)

// Ticket : description
type Ticket struct {
	gorm.Model
	TicketID      string  `gorm:"type:varchar(100)" json:"ticket_id"`
	PriorityID    int32   `json:"priority_id"`
	StatusID      int32   `json:"status_id"`
	CategoryID    int32   `json:"category_id"`
	SubCategoryID int32   `json:"sub_category_id"`
	Author        User    `json:"author"`
	UserID        uint32  `json:"user_id"`
	Holder        User    `json:"holder"`
	HolderID      uint32  `json:"holder_id"`
	Closer        User    `json:"closer"`
	CloserID      uint32  `json:"closer_id"`
	Assignee      User    `json:"assignee"`
	AssigneeID    uint32  `json:"assignee_id"`
	ContactNo     string  `json:"contact_no"`
	Title         string  `json:"title"`
	Message       string  `json:"message"`
	Justification string  `json:"justification"`
	Addresss      string  `json:"address"`
	Lat           float64 `json:"lat"`
	Lng           float64 `json:"lng"`
	Held          bool    `gorm:"default:false" json:"held"`
	Resolved      bool    `gorm:"default:false" json:"resolved"`
}

// HashTID : description
func HashTID(tid uint) string {
	hd := hashids.NewData()
	hd.Salt = os.Getenv("SALT")
	hd.MinLength = 6
	h, _ := hashids.NewWithData(hd)
	e, _ := h.Encode([]int{int(tid)})

	return e
}

// Prepare : description
func (t *Ticket) Prepare() {

	if t.ContactNo != "" {
		t.ContactNo = phonenumber.Parse(t.ContactNo, "KE")
	}
	t.Title = html.EscapeString(strings.TrimSpace(t.Title))
	t.Message = html.EscapeString(strings.TrimSpace(t.Message))
	t.Author = User{}
}

// Validate : description
func (t *Ticket) Validate() error {

	if t.Title == "" {
		return errors.New("Required Title")
	}
	if t.Message == "" {
		return errors.New("Required Message")
	}
	if t.ContactNo == "" {
		return errors.New("Required Contact No")
	}
	if t.UserID < 1 {
		return errors.New("Required Author")
	}
	return nil
}

// SaveTicket : description
func (t *Ticket) SaveTicket(db *gorm.DB) (*Ticket, error) {
	var err, uperr error
	err = db.Debug().Model(&Ticket{}).Create(&t).Error
	if err != nil {
		return &Ticket{}, err
	}
	if t.ID != 0 {
		if t.TicketID == "" {
			newHash := HashTID(t.ID)
			uperr = db.Debug().Model(&Ticket{}).Where("id = ?", t.ID).Updates(Ticket{TicketID: newHash}).Error

			if uperr != nil {
				fmt.Println(uperr)
			}
		}

		err = db.Debug().Model(&User{}).Select("ID, Email").Where("id = ?", t.UserID).Take(&t.Author).Error
		if err != nil {
			return &Ticket{}, err
		}
	}
	return t, nil
}

// FindAllTickets : description
func (t *Ticket) FindAllTickets(db *gorm.DB) (*[]Ticket, error) {
	var err error
	tickets := []Ticket{}
	err = db.Debug().Model(&Ticket{}).Limit(100).Find(&tickets).Error
	if err != nil {
		return &[]Ticket{}, err
	}
	if len(tickets) > 0 {
		for i := range tickets {
			err := db.Debug().Model(&User{}).Select("ID, Email").Where("id = ?", tickets[i].UserID).Take(&tickets[i].Author).Error
			if err != nil {
				return &[]Ticket{}, err
			}
		}
	}
	return &tickets, nil
}

// FindTicketByID : description
func (t *Ticket) FindTicketByID(db *gorm.DB, tid uint64) (*Ticket, error) {
	var err error
	err = db.Debug().Model(&Ticket{}).Where("id = ?", tid).Take(&t).Error
	if err != nil {
		return &Ticket{}, err
	}
	if t.ID != 0 {
		err = db.Debug().Model(&User{}).Select("ID, Email").Where("id = ?", t.UserID).Take(&t.Author).Error
		if err != nil {
			return &Ticket{}, err
		}
	}
	return t, nil
}

// UpdateATicket : description
func (t *Ticket) UpdateATicket(db *gorm.DB) (*Ticket, error) {

	var err error

	err = db.Debug().Model(&Ticket{}).Where("id = ?", t.ID).Updates(Ticket{
		CloserID:      t.CloserID,
		AssigneeID:    t.AssigneeID,
		HolderID:      t.HolderID,
		Justification: t.Justification,
		Held:          t.Held,
		Resolved:      t.Resolved,
	}).Error
	if err != nil {
		return &Ticket{}, err
	}

	if t.CloserID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.CloserID).Take(&t.Closer).Error
		if err != nil {
			return &Ticket{}, err
		}
	}

	if t.AssigneeID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.AssigneeID).Take(&t.Assignee).Error
		if err != nil {
			return &Ticket{}, err
		}
	}

	if t.HolderID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.CloserID).Take(&t.Holder).Error
		if err != nil {
			return &Ticket{}, err
		}
	}
	return t, nil
}
