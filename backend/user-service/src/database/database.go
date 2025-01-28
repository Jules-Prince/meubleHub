package database

import (
	"hexagone/user-service/src/models"
	"hexagone/user-service/src/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB


func ConnectDatabase(dbPath string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
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
