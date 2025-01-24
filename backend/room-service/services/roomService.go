package services

import (
	"hexagone/room-service/database"
	"hexagone/room-service/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateRoomInput struct {
	Name   string `json:"name" binding:"required"`
	HomeID uint   `json:"home_id" binding:"required"`
}

// CreateRoom handles the creation of a new room linked to a home
func CreateRoom(c *gin.Context) {
	var input CreateRoomInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the room
	room := models.Room{Name: input.Name, HomeID: input.HomeID}
	if result := database.DB.Create(&room); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": room})
}

// ListRooms handles fetching all rooms for a specific home
func ListRooms(c *gin.Context) {
	homeIDStr := c.DefaultQuery("home_id", "")
	if homeIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "home_id is required"})
		return
	}

	homeID, err := strconv.ParseUint(homeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid home_id"})
		return
	}

	var rooms []models.Room
	if err := database.DB.Where("home_id = ?", homeID).Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rooms"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": rooms})
}
