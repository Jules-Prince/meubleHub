package database

import (
	"hexagone/home-service/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("data/home.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Home database connected successfully!")

	// Migrate the schema
	err = DB.AutoMigrate(&models.Home{})
	if err != nil {
		log.Fatal("Failed to migrate home database:", err)
	}
}