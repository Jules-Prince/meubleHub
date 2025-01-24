package main

import (
	"hexagone/room-service/database"
	"hexagone/room-service/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	database.ConnectDatabase()

	r := gin.Default()

	// Routes
	r.POST("/rooms", services.CreateRoom)
	r.GET("/rooms", services.ListRooms) 

	// Start the service
	r.Run(":8080") 
}
