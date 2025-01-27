package main

import (
	"hexagone/room-service/database"
	"hexagone/room-service/services"
	"hexagone/room-service/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitLogger()
	utils.Log.Info("Starting Object Service")

	// Connect to the database
	database.ConnectDatabase()
	utils.Log.Info("Connected to SQLite")

	r := gin.Default()

	// Routes
	r.POST("/rooms", services.CreateRoom)
	r.GET("/rooms", services.ListRooms) 

	utils.Log.Info("Starting HTTP server on port 8080")
	// Start the service
	r.Run(":8080") 
}
