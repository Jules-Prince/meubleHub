package database

import (
	"hexagone/room-service/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("data/room.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Room database connected successfully!")

	// Migrate the schema for Room
	err = DB.AutoMigrate(&models.Room{})
	if err != nil {
		log.Fatal("Failed to migrate room database:", err)
	}
}
