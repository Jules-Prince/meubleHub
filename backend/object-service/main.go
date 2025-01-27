package main

import (
	"hexagone/object-service/database"
	"hexagone/object-service/services"
	"hexagone/object-service/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitLogger()
	utils.Log.Info("Starting Object Service")
	
	// Connect to DragonflyDB
	database.ConnectDatabase()
	utils.Log.Info("Connected to DragonflyDB")

	r := gin.Default()

	// Object routes
	r.POST("/objects", services.CreateObject)               // Add a new object
	r.GET("/objects", services.ListObjects)                 // List all objects
	r.PATCH("/objects/:id/reserve", services.ReserveObject) // Reserve an object
	r.GET("/objects/reserved", services.ListReservedObjects) // List all reserved objects


	// Start the service
	utils.Log.Info("Starting HTTP server on port 8083")
	r.Run(":8083") // Different port for the object service
}
