package services

import (
	"hexagone/home-service/database"
	"hexagone/home-service/models"
	"hexagone/home-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)
type CreateHomeInput struct {
	Name string `json:"name" binding:"required"`
}

// CreateHome handles the creation of a new home
func CreateHome(c *gin.Context) {
	var input CreateHomeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error binding JSON in CreateHome")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.Log.WithFields(logrus.Fields{
		"name": input.Name,
	}).Info("Creating home")

	// Create the home
	home := models.Home{Name: input.Name}
	if result := database.DB.Create(&home); result.Error != nil {
		utils.Log.WithFields(logrus.Fields{
			"name":  input.Name,
			"error": result.Error.Error(),
		}).Error("Error creating home in the database")
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	utils.Log.WithFields(logrus.Fields{
		"id":   home.ID,
		"name": home.Name,
	}).Info("Home created successfully")

	c.JSON(http.StatusOK, gin.H{"data": home})
}

// ListHomes handles fetching all homes
func ListHomes(c *gin.Context) {
	utils.Log.Info("Fetching all homes")

	var homes []models.Home
	if err := database.DB.Find(&homes).Error; err != nil {
		utils.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to retrieve homes from the database")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve homes"})
		return
	}

	utils.Log.WithFields(logrus.Fields{
		"count": len(homes),
	}).Info("Homes fetched successfully")

	c.JSON(http.StatusOK, gin.H{"data": homes})
}