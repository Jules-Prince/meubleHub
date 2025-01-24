package main

import (
	"hexagone/home-service/database"
	"hexagone/home-service/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	database.ConnectDatabase()

	r := gin.Default()

	// Routes
	r.POST("/homes", services.CreateHome)
	r.GET("/homes", services.ListHomes)

	// Start the server
	r.Run(":8080")
}