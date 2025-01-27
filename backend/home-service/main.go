package main

import (
	"fmt"
	"hexagone/home-service/database"
	"hexagone/home-service/services"
	"hexagone/home-service/utils"
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

	// Initialize logger
	utils.InitLogger()

	// Log the loaded configurations
	utils.Log.Infof("Starting Home Service on port %s with DB at %s", port, dbPath)

	// Initialize database
	database.ConnectDatabase(dbPath)

	// Set up Gin router
	r := gin.Default()
	utils.Log.Info("Starting Home Service")

	// Routes
	r.POST("/homes", services.CreateHome)
	r.GET("/homes", services.ListHomes)

	utils.Log.Info("Starting HTTP server on port 8080")


	// Start the server
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		utils.Log.Fatalf("Failed to start server: %v", err)
	}
}