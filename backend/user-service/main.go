package main

import (
	"hexagone/user-service/database"
	"hexagone/user-service/services"
	"hexagone/user-service/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitLogger()
	utils.Log.Info("Starting User service")

	// Connect to the database
	database.ConnectDatabase()
	utils.Log.Info("Connected to SQLite")

	r := gin.Default()

	// Routes
	r.POST("/users", services.CreateUser)  // Create a user
	r.POST("/login", services.Login)       // Login
	r.GET("/users", services.ListUsers)    // List all users

	utils.Log.Info("Starting HTTP server on port 8080")
	// Start the server
	r.Run(":8080")
}
