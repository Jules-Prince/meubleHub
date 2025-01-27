package main

import (
	"fmt"
	"hexagone/object-service/database"
	"hexagone/object-service/services"
	"hexagone/object-service/utils"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	
	port := os.Getenv("PORT")
	if port == "" {
		utils.Log.Error("PORT is not set in the environment variables")
	}

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
	utils.Log.Infof("Starting HTTP server on port %s", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		utils.Log.Fatalf("Failed to start server: %v", err)
	}
}
