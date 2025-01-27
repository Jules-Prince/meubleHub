package main

import (
	"hexagone/object-service/database"
	"hexagone/object-service/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to DragonflyDB
	database.ConnectDatabase()

	r := gin.Default()

	// Object routes
	r.POST("/objects", services.CreateObject)               // Add a new object
	r.GET("/objects", services.ListObjects)                 // List all objects
	r.PATCH("/objects/:id/reserve", services.ReserveObject) // Reserve an object
	r.GET("/objects/reserved", services.ListReservedObjects) // List all reserved objects


	// Start the service
	r.Run(":8083") // Different port for the object service
}
