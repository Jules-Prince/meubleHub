package main

import (
	"fmt"
	"hexagone/room-service/database"
	"hexagone/room-service/services"
	"hexagone/room-service/utils"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		utils.Log.Error("Warning: No .env file found")
	}

	// Get configuration from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		utils.Log.Error("PORT is not set in the environment variables")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		utils.Log.Error("DB_PATH is not set in the environment variables")
	}
	
	utils.InitLogger()
	utils.Log.Info("Starting Object Service")

	// Connect to the database
	database.ConnectDatabase(dbPath)
	utils.Log.Info("Connected to SQLite")

	r := gin.Default()

	// Routes
	r.POST("/rooms", services.CreateRoom)
	r.GET("/rooms", services.ListRooms) 

	utils.Log.Infof("Starting HTTP server on port %s", port)

	// Start the server
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		utils.Log.Fatalf("Failed to start server: %v", err)
	}
}
