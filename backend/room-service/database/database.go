package database

import (
	"hexagone/room-service/models"
	"hexagone/room-service/utils"
	
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("data/room.db"), &gorm.Config{})
	if err != nil {
		utils.Log.WithField("error", err.Error()).Error("Failed to connect to database:")
	}

	utils.Log.Info("Room database connected successfully!")

	// Migrate the schema for Room
	err = DB.AutoMigrate(&models.Room{})
	if err != nil {
		utils.Log.WithField("error", err.Error()).Error("Failed to connect to database:")
	}
}
