package jobs

import (
	"fmt"
	"time"

	"github.com/Tchayo/gotickets/api/models"
	"github.com/jinzhu/gorm"
)

// UpdateSLA : function to update ticket priority time
func UpdateSLA(db *gorm.DB) {
	var err error
	timenow := time.Now().Local()

	fmt.Printf("Time now is %s\r\n", timenow)

	futuretime := timenow.Add(time.Hour*time.Duration(2) + time.Minute*time.Duration(20))

	fmt.Printf("Future time is %s\r\n", futuretime)

	tickets := []models.Ticket{}
	err = db.Debug().Preload("SubCategory.Category").Limit(100).Find(&tickets).Error
	if err != nil {
		println(err)
	}

	for i := range tickets {
		println(tickets[i].SubCategory.Category.TeamID)
	}

}
