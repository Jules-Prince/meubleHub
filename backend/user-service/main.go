package main

import (
	"fmt"
	"hexagone/user-service/database"
	"hexagone/user-service/services"
	"hexagone/user-service/utils"
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
	
	utils.InitLogger()
	utils.Log.Info("Starting User service")

	// Connect to the database
	database.ConnectDatabase(dbPath)
	utils.Log.Info("Connected to SQLite")

	r := gin.Default()

	// Routes
	r.POST("/users", services.CreateUser)  // Create a user
	r.POST("/login", services.Login)       // Login
	r.GET("/users", services.ListUsers)    // List all users

	utils.Log.Infof("Starting HTTP server on port %s", port)

	// Start the server
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		utils.Log.Fatalf("Failed to start server: %v", err)
	}
}
