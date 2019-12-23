package seed

import (
	"log"

	"github.com/Tchayo/gotickets/api/models"
	"github.com/jinzhu/gorm"
)

var teams = []models.Team{
	models.Team{
		Name:    "Support",
		Company: "Adtel",
		Email:   "felix.achayo@adtel.co.ke",
	},
	models.Team{
		Name:    "Finance",
		Company: "Adtel",
		Email:   "felix.achayo@adtel.co.ke",
	},
}

var users = []models.User{
	models.User{
		Email:    "felix.achayo@adtel.co.ke",
		Password: "PassWord*",
	},
	models.User{
		Email:    "achayof@gmail.com",
		Password: "PassWord*",
	},
}

var priorites = []models.Priority{
	models.Priority{
		Title:          "Normal",
		TimeoutHours:   1,
		TimeoutMinutes: 0,
		Color:          "#ccc",
		Description:    "Normal priority",
		Fin:            false,
	},
	models.Priority{
		Title:          "Overdue",
		TimeoutHours:   2,
		TimeoutMinutes: 0,
		Color:          "#FFC107",
		Description:    "Overdue priority",
		Fin:            false,
	},
	models.Priority{
		Title:          "Critical",
		TimeoutHours:   3,
		TimeoutMinutes: 0,
		Color:          "#E91E63",
		Description:    "Critical priority",
		Fin:            false,
	},
	models.Priority{
		Title:          "Resolved",
		TimeoutHours:   0,
		TimeoutMinutes: 0,
		Color:          "#4CAF50",
		Description:    "Resolved priority",
		Fin:            true,
	},
}

var statuses = []models.Status{
	models.Status{
		Title:       "Normal",
		Color:       "#BDBDBD",
		Description: "Normal ticket status",
		Resolution:  false,
		Lockable:    false,
		Assignable:  true,
	},
	models.Status{
		Title:       "Re-assign",
		Color:       "#03A9F4",
		Description: "Re-assign ticket status",
		Resolution:  false,
		Lockable:    false,
		Assignable:  true,
	},
	models.Status{
		Title:       "Lock",
		Color:       "#3F51B5",
		Description: "Lock ticket status",
		Resolution:  false,
		Lockable:    true,
		Assignable:  false,
	},
	models.Status{
		Title:       "Resolved",
		Color:       "#4CAF50",
		Description: "Resolved ticket status",
		Resolution:  true,
		Lockable:    false,
		Assignable:  true,
	},
}

var categories = []models.Category{
	models.Category{
		TeamID:      1,
		Title:       "Repost",
		Description: "Repost related issues",
	},
	models.Category{
		TeamID:      1,
		Title:       "CRM",
		Description: "CRM related issues",
	},
	models.Category{
		TeamID:      2,
		Title:       "Reversal",
		Description: "Reversal related issues",
	},
}

var subs = []models.Sub{
	models.Sub{
		CategoryID:  1,
		Title:       "Mpesa Repost",
		Description: "Repost related issues",
	},
	models.Sub{
		CategoryID:  2,
		Title:       "CRM Timeout",
		Description: "CRM timeout related issues",
	},
	models.Sub{
		CategoryID:  3,
		Title:       "Mpesa Reversal",
		Description: "Mpesa reversal related issues",
	},
}

// Load : load migrations
func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Sub{}, &models.Category{}, &models.Priority{}, &models.User{}, &models.Team{}, &models.Status{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	err = db.Debug().AutoMigrate(&models.Team{}, &models.User{}, &models.Priority{}, &models.Status{}, &models.Category{}, &models.Sub{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.User{}).AddForeignKey("team_id", "teams(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	err = db.Debug().Model(&models.Priority{}).AddForeignKey("team_id", "teams(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	err = db.Debug().Model(&models.Category{}).AddForeignKey("team_id", "teams(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	err = db.Debug().Model(&models.Sub{}).AddForeignKey("category_id", "categories(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	var teamSelect uint

	for i := range teams {
		err = db.Debug().Model(&models.Team{}).Create(&teams[i]).Error
		if err != nil {
			log.Fatalf("cannot seed teams table: %v", err)
		}
		users[i].TeamID = teams[i].ID
		teamSelect = teams[0].ID

		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}

	}

	for i := range priorites {
		priorites[i].TeamID = teamSelect

		err = db.Debug().Model(&models.Priority{}).Create(&priorites[i]).Error
		if err != nil {
			log.Fatalf("cannot seed priorities table: %v", err)
		}
	}

	for i := range statuses {
		err = db.Debug().Model(&models.Status{}).Create(&statuses[i]).Error
		if err != nil {
			log.Fatalf("cannot seed statuses table: %v", err)
		}
	}

	for i := range categories {
		err = db.Debug().Model(&models.Category{}).Create(&categories[i]).Error
		if err != nil {
			log.Fatalf("cannot seed categories table: %v", err)
		}
	}

	for i := range subs {
		err = db.Debug().Model(&models.Sub{}).Create(&subs[i]).Error
		if err != nil {
			log.Fatalf("cannot seed sub categories table: %v", err)
		}
	}
}
