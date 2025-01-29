package main

import (
	"fmt"
	"hexagone/room-service/src/database"
	"hexagone/room-service/src/services"
	"hexagone/room-service/src/utils"
	"os"

	"github.com/gin-contrib/cors"
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
	utils.Log.Info("Starting Object Service")

	// Connect to the database
	database.ConnectDatabase(dbPath)
	utils.Log.Info("Connected to SQLite")

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Routes
	r.POST("/rooms", services.CreateRoom)
	r.GET("/rooms", services.ListRooms) 

	utils.Log.Infof("Starting HTTP server on port %s", port)

	// Start the server
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		utils.Log.Fatalf("Failed to start server: %v", err)
	}
}
