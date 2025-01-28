package main

import (
	"fmt"
	"hexagone/home-service/src/database"
	"hexagone/home-service/src/services"
	"hexagone/home-service/src/utils"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
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
	
	// Initialize database
	database.ConnectDatabase(dbPath)

	// Set up Gin router
	r := gin.Default()
	utils.Log.Info("Starting Home Service")

	// Routes
	r.POST("/homes", services.CreateHome)
	r.GET("/homes", services.ListHomes)

	utils.Log.Infof("Starting HTTP server on port %s", port)

	// Start the server
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		utils.Log.Fatalf("Failed to start server: %v", err)
	}
}