package main

import (
	"hexagone/user-service/database"
	"hexagone/user-service/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	database.ConnectDatabase()

	r := gin.Default()

	// Routes
	r.POST("/users", services.CreateUser)  // Create a user
	r.POST("/login", services.Login)       // Login
	r.GET("/users", services.ListUsers)    // List all users


	// Start the server
	r.Run(":8080")
}
