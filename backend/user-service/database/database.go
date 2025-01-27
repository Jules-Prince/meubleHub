package database

import (
	"hexagone/user-service/models"
	"hexagone/user-service/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB


func ConnectDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("data/usrer.db"), &gorm.Config{})
	if err != nil {
		utils.Log.WithField("error", err.Error()).Error("Failed to connect to database")
	}

	utils.Log.Info("Room database connected successfully!")

	// Migrate the schema for Room
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		utils.Log.WithField("error", err.Error()).Error("Failed to connect to database")
	}
}
