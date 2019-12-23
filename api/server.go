package api

import (
	"fmt"
	"log"
	"os"

	"github.com/Tchayo/gotickets/api/controllers"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}

// Run : description
func Run() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("Getting env values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	// uncomment to seed sample data
	// seed.Load(server.DB)

	// test chron function
	// jobs.UpdateSLA(server.DB)

	server.Run(":8080")
}
