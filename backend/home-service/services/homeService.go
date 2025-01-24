package services

import (
	"hexagone/home-service/database"
	"hexagone/home-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateHomeInput struct {
	Name string `json:"name" binding:"required"`
}

// CreateHome handles the creation of a new home
func CreateHome(c *gin.Context) {
	var input CreateHomeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the home
	home := models.Home{Name: input.Name}
	if result := database.DB.Create(&home); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": home})
}

func ListHomes(c *gin.Context) {
	var homes []models.Home
	if err := database.DB.Find(&homes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve homes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": homes})
}
