package models

import (
	"errors"
	"fmt"
	"html"
	"log"
	"os"
	"strings"

	"github.com/Tchayo/gotickets/api/utils/filter"
	"github.com/Tchayo/gotickets/api/utils/formatdate"
	"github.com/Tchayo/gotickets/api/utils/sendmail"
	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/dongri/phonenumber"
	"github.com/jinzhu/gorm"
	"github.com/speps/go-hashids"
)

// Ticket : description
type Ticket struct {
	gorm.Model
	TicketID string `gorm:"type:varchar(100)" json:"ticket_id"`
	// Priority   Priority `gorm:"auto_preload" json:"priority"`
	PriorityID int32 `json:"priority_id"`
	StatusID   int32 `json:"status_id"`
	// Category      Category `json:"category"`
	// CategoryID    int32    `json:"category_id"`
	SubCategory   Sub     `gorm:"auto_preload" json:"sub_category"`
	SubCategoryID int32   `json:"sub_category_id"`
	User          User    `gorm:"auto_preload"  json:"user"`
	UserID        uint32  `json:"user_id"`
	Holder        User    `gorm:"auto_preload" json:"holder"`
	HolderID      uint32  `json:"holder_id"`
	Closer        User    `gorm:"auto_preload" json:"closer"`
	CloserID      uint32  `json:"closer_id"`
	Assignee      User    `gorm:"auto_preload" json:"assignee"`
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
	t.User = User{}
	// t.Priority = Priority{}
	t.SubCategory = Sub{}
}

// Validate : description
func (t *Ticket) Validate(action string) error {

	switch strings.ToLower(action) {
	case "update":
		if t.Justification == "" {
			return errors.New("Required Justification")
		}
		return nil
	default:
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
			return errors.New("Required User")
		}
		if t.SubCategoryID < 1 {
			return errors.New("Required Sub-category")
		}
		return nil
	}
}

// SaveTicket : description
func (t *Ticket) SaveTicket(db *gorm.DB) (*Ticket, error) {
	var err, updaterr error

	if tusr := db.Where("id = ?", t.UserID).First(&User{}); tusr.Error != nil {
		return &Ticket{}, tusr.Error
	}

	if tsub := db.Where("id = ?", t.SubCategoryID).First(&Sub{}); tsub.Error != nil {
		return &Ticket{}, tsub.Error
	}

	err = db.Debug().Model(&Ticket{}).Create(&t).Error
	if err != nil {
		return &Ticket{}, err
	}
	if t.ID != 0 {
		if t.TicketID == "" {
			newHash := HashTID(t.ID)
			updaterr = db.Debug().Model(&Ticket{}).Where("id = ?", t.ID).Updates(Ticket{TicketID: newHash}).Error

			if updaterr != nil {
				fmt.Println(updaterr)
			}
		}

		err = db.Debug().Where("id = ?", t.ID).Preload("User").Preload("User.Team").Preload("Assignee").Preload("Holder").Preload("Closer").Preload("SubCategory").Preload("SubCategory.Category").Preload("SubCategory.Category.Team").Find(&t).Error
		if err != nil {
			return &Ticket{}, err
		}

		// send mail to user --- later modified to team mail
		if usermail := t.User.Email; usermail != "" {
			formattedDate := formatdate.FormatDate(t.CreatedAt, "RFCN")
			mailUser := t.User.Email

			m := sendmail.Mail{ToAddr: usermail,
				FromName: "Ticketing System",
				FromAddr: "felix.achayo@adtel.co.ke",
				Subject:  "Ticket Created",
				Body:     fmt.Sprintf("Ticket subject: %s. \r\n\n%s. \r\n\nCreated @ %s by %s", t.Title, t.Message, formattedDate, mailUser)}

			err := sendmail.Mailer(m)
			if err != nil {
				log.Fatal(err)
			}
		}

	}
	return t, nil
}

// FindAllTickets : description
func (t *Ticket) FindAllTickets(db *gorm.DB, f *filter.Filter) (*pagination.Paginator, error) {

	var page, limit int
	tickets := []Ticket{}

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
		DB:      db.Preload("User").Preload("User.Team").Preload("Assignee").Preload("Holder").Preload("Closer").Preload("SubCategory").Preload("SubCategory.Category").Preload("SubCategory.Category.Team"),
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"id desc"},
	}, &tickets)

	return res, nil

}

// FindTicketByID : description
func (t *Ticket) FindTicketByID(db *gorm.DB, tid uint64) (*Ticket, error) {
	var err error
	err = db.Debug().Where("id = ?", tid).Preload("User").Preload("User.Team").Preload("Assignee").Preload("Holder").Preload("Closer").Preload("SubCategory").Preload("SubCategory.Category").Preload("SubCategory.Category.Team").Find(&t).Error
	if err != nil {
		return &Ticket{}, err
	}
	return t, nil
}

// UpdateATicket : description
func (t *Ticket) UpdateATicket(db *gorm.DB) (*Ticket, error) {

	var err error

	if tusr := db.Where("id = ? AND resolved = ?", t.ID, false).First(&Ticket{}); tusr.Error != nil {
		return &Ticket{}, errors.New("Ticket already resolved, create a new one instead")
	}

	assigneeOld := t.Assignee.Email

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

	err = db.Debug().Where("id = ?", t.ID).Preload("User").Preload("User.Team").Preload("Assignee").Preload("Holder").Preload("Closer").Preload("SubCategory").Preload("SubCategory.Category").Preload("SubCategory.Category.Team").Find(&t).Error
	if err != nil {
		return &Ticket{}, err
	}

	if t.Resolved == true {
		// send mail to user --- later modified to team mail
		if usermail := t.Closer.Email; usermail != "" {
			formattedDate := formatdate.FormatDate(t.UpdatedAt, "RFC")
			mailUser := t.Closer.Email

			m := sendmail.Mail{ToAddr: usermail,
				FromName: "Ticketing System",
				FromAddr: "felix.achayo@adtel.co.ke",
				Subject:  "Ticket Closed",
				Body:     fmt.Sprintf("Ticket subject: %s. \r\n\n%s. \r\n\nClosed @ %s by %s", t.Title, t.Message, formattedDate, mailUser)}

			err := sendmail.Mailer(m)
			if err != nil {
				log.Fatal(err)
			}
		}

	}

	if assigneeOld != t.Assignee.Email && t.Assignee.Email != "" {
		// send mail to user --- later modified to team mail
		if usermail := t.Assignee.Email; usermail != "" {
			formattedDate := formatdate.FormatDate(t.UpdatedAt, "RFCN")
			mailUser := t.Assignee.Email

			m := sendmail.Mail{ToAddr: usermail,
				FromName: "Ticketing System",
				FromAddr: "felix.achayo@adtel.co.ke",
				Subject:  "Ticket assigned to you",
				Body:     fmt.Sprintf("Ticket subject: %s. \r\n\n%s. \r\n\nCreated @ %s by %s", t.Title, t.Message, formattedDate, mailUser)}

			err := sendmail.Mailer(m)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if t.HolderID > 0 {
		// send mail to user --- later modified to team mail
		// if usermail := t.Holder.Email; usermail != "" {
		// 	formattedDate := formatdate.FormatDate(t.CreatedAt, "RFCN")
		// 	mailUser := t.Holder.Email

		// 	m := sendmail.Mail{ToAddr: usermail,
		// 		FromName: "Ticketing System",
		// 		FromAddr: "felix.achayo@adtel.co.ke",
		// 		Subject:  "Ticket Created",
		// 		Body:     fmt.Sprintf("Ticket subject: %s. \r\n\n%s. \r\n\nCreated @ %s by %s", t.Title, t.Message, formattedDate, mailUser)}

		// 	err := sendmail.Mailer(m)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// }
	}
	return t, nil
}
