package services

import (
	"hexagone/room-service/src/database"
	"hexagone/room-service/src/models"
	"hexagone/room-service/src/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)
type CreateRoomInput struct {
	Name   string `json:"name" binding:"required"`
	HomeID uint   `json:"home_id" binding:"required"`
}

// CreateRoom handles the creation of a new room linked to a home
func CreateRoom(c *gin.Context) {
	var input CreateRoomInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error binding JSON in CreateRoom")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.Log.WithFields(logrus.Fields{
		"name":   input.Name,
		"homeID": input.HomeID,
	}).Info("Creating room")

	// Create the room
	room := models.Room{Name: input.Name, HomeID: input.HomeID}
	if result := database.DB.Create(&room); result.Error != nil {
		utils.Log.WithFields(logrus.Fields{
			"error": result.Error.Error(),
		}).Error("Error creating room in the database")
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	utils.Log.WithFields(logrus.Fields{
		"id":     room.ID,
		"name":   room.Name,
		"homeID": room.HomeID,
	}).Info("Room created successfully")

	c.JSON(http.StatusOK, gin.H{"data": room})
}

// ListRooms handles fetching all rooms for a specific home
func ListRooms(c *gin.Context) {
	homeIDStr := c.DefaultQuery("home_id", "")
	if homeIDStr == "" {
		utils.Log.Warn("home_id is missing in ListRooms request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "home_id is required"})
		return
	}

	homeID, err := strconv.ParseUint(homeIDStr, 10, 32)
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"homeIDStr": homeIDStr,
			"error":     err.Error(),
		}).Error("Invalid home_id in ListRooms request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid home_id"})
		return
	}

	utils.Log.WithFields(logrus.Fields{
		"homeID": homeID,
	}).Info("Fetching rooms for home")

	var rooms []models.Room
	if err := database.DB.Where("home_id = ?", homeID).Find(&rooms).Error; err != nil {
		utils.Log.WithFields(logrus.Fields{
			"homeID": homeID,
			"error":  err.Error(),
		}).Error("Failed to retrieve rooms from the database")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rooms"})
		return
	}

	utils.Log.WithFields(logrus.Fields{
		"homeID": homeID,
		"count":  len(rooms),
	}).Info("Rooms fetched successfully")

	c.JSON(http.StatusOK, gin.H{"data": rooms})
}