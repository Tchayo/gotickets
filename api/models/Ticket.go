package models

import (
	"errors"
	"html"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/speps/go-hashids"
)

type Ticket struct {
	gorm.Model
	TicketID      string  `gorm:"type:varchar(100);unique_index" json:"ticket_id"`
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
	Held          bool    `json:"held"`
	Resolved      bool    `json:"resolved"`
}

func HashTID(tid uint) string {
	hd := hashids.NewData()
	hd.Salt = os.Getenv("SALT")
	hd.MinLength = 6
	h, _ := hashids.NewWithData(hd)
	e, _ := h.Encode([]int{int(tid)})

	return e
}

func (t *Ticket) Prepare() {
	t.TicketID = HashTID(t.ID)
	t.Title = html.EscapeString(strings.TrimSpace(t.Title))
	t.Message = html.EscapeString(strings.TrimSpace(t.Message))
	t.Author = User{}
}

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

func (t *Ticket) SaveTicket(db *gorm.DB) (*Ticket, error) {
	var err error
	err = db.Debug().Model(&Ticket{}).Create(&t).Error
	if err != nil {
		return &Ticket{}, err
	}
	if t.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.UserID).Take(&t.Author).Error
		if err != nil {
			return &Ticket{}, err
		}
	}
	return t, nil
}

func (t *Ticket) FindAllTickets(db *gorm.DB) (*[]Ticket, error) {
	var err error
	tickets := []Ticket{}
	err = db.Debug().Model(&Ticket{}).Limit(100).Find(&tickets).Error
	if err != nil {
		return &[]Ticket{}, err
	}
	if len(tickets) > 0 {
		for i := range tickets {
			err := db.Debug().Model(&User{}).Where("id = ?", tickets[i].UserID).Take(&tickets[i].Author).Error
			if err != nil {
				return &[]Ticket{}, err
			}
		}
	}
	return &tickets, nil
}

func (t *Ticket) FindTicketByID(db *gorm.DB, tid uint64) (*Ticket, error) {
	var err error
	err = db.Debug().Model(&Ticket{}).Where("id = ?", tid).Take(&t).Error
	if err != nil {
		return &Ticket{}, err
	}
	if t.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.UserID).Take(&t.Author).Error
		if err != nil {
			return &Ticket{}, err
		}
	}
	return t, nil
}

func (t *Ticket) CloseATicket(db *gorm.DB) (*Ticket, error) {

	var err error

	err = db.Debug().Model(&Ticket{}).Where("id = ?", t.ID).Updates(Ticket{CloserID: t.CloserID, Justification: t.Justification, Resolved: true}).Error
	if err != nil {
		return &Ticket{}, err
	}
	if t.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.CloserID).Take(&t.Closer).Error
		if err != nil {
			return &Ticket{}, err
		}
	}
	return t, nil
}
