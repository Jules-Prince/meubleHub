package main

import (
	"hexagone/home-service/database"
	"hexagone/home-service/services"
	"hexagone/home-service/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitLogger()
	utils.Log.Info("Starting Home Service")

	// Connect to the database
	database.ConnectDatabase()
	utils.Log.Info("Connected to SQLite")

	r := gin.Default()

	// Routes
	r.POST("/homes", services.CreateHome)
	r.GET("/homes", services.ListHomes)

	utils.Log.Info("Starting HTTP server on port 8080")
	// Start the server
	r.Run(":8080")
}