package services

import (
	"encoding/json"
	"hexagone/object-service/database"
	"hexagone/object-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReserveObjectInput struct {
	UserID string `json:"userId" binding:"required"`
}

// ReserveObject allows a user to reserve an object by its ID
func ReserveObject(c *gin.Context) {
	objectID := c.Param("id")
	var input ReserveObjectInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the object
	val, err := database.RDB.Get(database.Ctx, objectID).Result()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Object not found"})
		return
	}

	var object models.Object
	if err := json.Unmarshal([]byte(val), &object); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process object data"})
		return
	}

	// Check if the object is already reserved
	if object.IsReserved {
		c.JSON(http.StatusConflict, gin.H{"error": "Object is already reserved"})
		return
	}

	// Reserve the object
	object.IsReserved = true
	object.ReservedBy = input.UserID

	// Save the updated object back to DragonflyDB
	data, err := json.Marshal(object)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save reservation"})
		return
	}

	err = database.RDB.Set(database.Ctx, object.ID, data, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update object in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": object})
}

// ListReservedObjects retrieves all reserved objects
func ListReservedObjects(c *gin.Context) {
	keys, err := database.RDB.Keys(database.Ctx, "*").Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch keys"})
		return
	}

	reservedObjects := []models.Object{}
	for _, key := range keys {
		val, err := database.RDB.Get(database.Ctx, key).Result()
		if err != nil {
			continue // Skip if object retrieval fails
		}

		var obj models.Object
		if err := json.Unmarshal([]byte(val), &obj); err != nil {
			continue // Skip if JSON unmarshalling fails
		}

		if obj.IsReserved {
			reservedObjects = append(reservedObjects, obj)
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": reservedObjects})
}

type CreateObjectInput struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required"`
}

// CreateObject adds a new object to DragonflyDB
func CreateObject(c *gin.Context) {
	var input CreateObjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create an object
	object := models.Object{
		ID:   uuid.New().String(),
		Name: input.Name,
		Type: input.Type,
	}

	// Serialize object to JSON
	data, err := json.Marshal(object)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process object"})
		return
	}

	// Store object in DragonflyDB
	err = database.RDB.Set(database.Ctx, object.ID, data, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store object in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": object})
}

// ListObjects retrieves all objects from DragonflyDB
func ListObjects(c *gin.Context) {
	keys, err := database.RDB.Keys(database.Ctx, "*").Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch keys"})
		return
	}

	objects := []models.Object{}
	for _, key := range keys {
		val, err := database.RDB.Get(database.Ctx, key).Result()
		if err != nil {
			continue // Skip if object retrieval fails
		}

		var obj models.Object
		if err := json.Unmarshal([]byte(val), &obj); err != nil {
			continue // Skip if JSON unmarshalling fails
		}

		objects = append(objects, obj)
	}

	c.JSON(http.StatusOK, gin.H{"data": objects})
}
