package services

import (
	"encoding/json"
	"hexagone/object-service/src/database"
	"hexagone/object-service/src/models"
	"hexagone/object-service/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ReserveObjectInput struct {
	UserID string `json:"userId" binding:"required"`
}

// ReserveObject allows a user to reserve an object by its ID
func ReserveObject(c *gin.Context) {
	objectID := c.Param("id")
	var input ReserveObjectInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.WithFields(logrus.Fields{
			"objectID": objectID,
			"error":    err.Error(),
		}).Error("Failed to bind input for reservation")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.Log.WithField("objectID", objectID).Info("Attempting to reserve object")

	// Fetch the object
	val, err := database.RDB.Get(database.Ctx, objectID).Result()
	if err != nil {
		utils.Log.WithField("objectID", objectID).Warn("Object not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Object not found"})
		return
	}

	var object models.Object
	if err := json.Unmarshal([]byte(val), &object); err != nil {
		utils.Log.WithField("objectID", objectID).Error("Failed to unmarshal object data")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process object data"})
		return
	}

	// Check if the object is already reserved
	if object.IsReserved {
		utils.Log.WithField("objectID", objectID).Info("Object is already reserved")
		c.JSON(http.StatusConflict, gin.H{"error": "Object is already reserved"})
		return
	}

	// Reserve the object
	object.IsReserved = true
	object.ReservedBy = input.UserID

	// Save the updated object back to DragonflyDB
	data, err := json.Marshal(object)
	if err != nil {
		utils.Log.WithField("objectID", objectID).Error("Failed to marshal updated object")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save reservation"})
		return
	}

	err = database.RDB.Set(database.Ctx, object.ID, data, 0).Err()
	if err != nil {
		utils.Log.WithField("objectID", objectID).Error("Failed to update object in database")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update object in database"})
		return
	}

	utils.Log.WithFields(logrus.Fields{
		"objectID": object.ID,
		"userID":   input.UserID,
	}).Info("Object reserved successfully")

	c.JSON(http.StatusOK, gin.H{"data": object})
}

// ListReservedObjects retrieves all reserved objects
func ListReservedObjects(c *gin.Context) {
	utils.Log.Info("Fetching all reserved objects")

	keys, err := database.RDB.Keys(database.Ctx, "*").Result()
	if err != nil {
		utils.Log.WithField("error", err.Error()).Error("Failed to fetch keys from DragonflyDB")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch keys"})
		return
	}

	utils.Log.WithField("keysCount", len(keys)).Info("Fetched keys from database")

	reservedObjects := []models.Object{}
	for _, key := range keys {
		val, err := database.RDB.Get(database.Ctx, key).Result()
		if err != nil {
			utils.Log.WithField("key", key).Warn("Failed to retrieve object from DragonflyDB, skipping")
			continue
		}

		var obj models.Object
		if err := json.Unmarshal([]byte(val), &obj); err != nil {
			utils.Log.WithField("key", key).Warn("Failed to unmarshal object data, skipping")
			continue
		}

		if obj.IsReserved {
			utils.Log.WithField("objectID", obj.ID).Info("Reserved object found")
			reservedObjects = append(reservedObjects, obj)
		}
	}

	utils.Log.WithField("reservedCount", len(reservedObjects)).Info("Reserved objects fetched successfully")
	c.JSON(http.StatusOK, gin.H{"data": reservedObjects})
}

type CreateObjectInput struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required"`
	RoomID string `json:"room_id" binding:"required"`
}

// CreateObject adds a new object to DragonflyDB
func CreateObject(c *gin.Context) {
	var input CreateObjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.WithField("error", err.Error()).Error("Failed to bind JSON input for object creation")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.Log.WithFields(logrus.Fields{
		"objectName": input.Name,
		"objectType": input.Type,
	}).Info("Creating new object")

	// Create an object
	object := models.Object{
        ID:     uuid.New().String(),
        Name:   input.Name,
        Type:   input.Type,
        RoomID: input.RoomID,
    }

	// Serialize object to JSON
	data, err := json.Marshal(object)
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"objectID": object.ID,
			"error":    err.Error(),
		}).Error("Failed to marshal object to JSON")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process object"})
		return
	}

	// Store object in DragonflyDB
	err = database.RDB.Set(database.Ctx, object.ID, data, 0).Err()
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"objectID": object.ID,
			"error":    err.Error(),
		}).Error("Failed to store object in DragonflyDB")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store object in database"})
		return
	}

	utils.Log.WithField("objectID", object.ID).Info("Object created and stored successfully")
	c.JSON(http.StatusOK, gin.H{"data": object})
}

// ListObjects retrieves all objects from DragonflyDB
func ListObjects(c *gin.Context) {
	utils.Log.Info("Fetching all objects from DragonflyDB")

	keys, err := database.RDB.Keys(database.Ctx, "*").Result()
	if err != nil {
		utils.Log.WithField("error", err.Error()).Error("Failed to fetch keys from DragonflyDB")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch keys"})
		return
	}

	utils.Log.WithField("keysCount", len(keys)).Info("Fetched keys from database")

	objects := []models.Object{}
	for _, key := range keys {
		val, err := database.RDB.Get(database.Ctx, key).Result()
		if err != nil {
			utils.Log.WithField("key", key).Warn("Failed to retrieve object from DragonflyDB, skipping")
			continue
		}

		var obj models.Object
		if err := json.Unmarshal([]byte(val), &obj); err != nil {
			utils.Log.WithField("key", key).Warn("Failed to unmarshal object data, skipping")
			continue
		}

		utils.Log.WithField("objectID", obj.ID).Info("Object fetched successfully")
		objects = append(objects, obj)
	}

	utils.Log.WithField("objectsCount", len(objects)).Info("Objects fetched successfully")
	c.JSON(http.StatusOK, gin.H{"data": objects})
}

func ListObjectsByRoom(c *gin.Context) {
    roomID := c.Query("room_id")
    if roomID == "" {
        utils.Log.Warn("room_id is missing in ListObjectsByRoom request")
        c.JSON(http.StatusBadRequest, gin.H{"error": "room_id is required"})
        return
    }

    utils.Log.WithField("roomID", roomID).Info("Fetching objects for room")

    keys, err := database.RDB.Keys(database.Ctx, "*").Result()
    if err != nil {
        utils.Log.WithField("error", err.Error()).Error("Failed to fetch keys from DragonflyDB")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch keys"})
        return
    }

    utils.Log.WithField("keysCount", len(keys)).Info("Fetched keys from database")

    objects := []models.Object{}
    for _, key := range keys {
        val, err := database.RDB.Get(database.Ctx, key).Result()
        if err != nil {
            utils.Log.WithField("key", key).Warn("Failed to retrieve object from DragonflyDB, skipping")
            continue
        }

        var obj models.Object
        if err := json.Unmarshal([]byte(val), &obj); err != nil {
            utils.Log.WithField("key", key).Warn("Failed to unmarshal object data, skipping")
            continue
        }

        // Only include objects that belong to the specified room
        if obj.RoomID == roomID {
            utils.Log.WithField("objectID", obj.ID).Info("Object in room found")
            objects = append(objects, obj)
        }
    }

    utils.Log.WithFields(logrus.Fields{
        "roomID": roomID,
        "count":  len(objects),
    }).Info("Room objects fetched successfully")

    c.JSON(http.StatusOK, gin.H{"data": objects})
}



func UnreserveObject(c *gin.Context) {
    objectID := c.Param("id")

    utils.Log.WithField("objectID", objectID).Info("Attempting to unreserve object")

    // Fetch the object
    val, err := database.RDB.Get(database.Ctx, objectID).Result()
    if err != nil {
        utils.Log.WithField("objectID", objectID).Warn("Object not found")
        c.JSON(http.StatusNotFound, gin.H{"error": "Object not found"})
        return
    }

    var object models.Object
    if err := json.Unmarshal([]byte(val), &object); err != nil {
        utils.Log.WithField("objectID", objectID).Error("Failed to unmarshal object data")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process object data"})
        return
    }

    // Check if the object is not reserved
    if !object.IsReserved {
        utils.Log.WithField("objectID", objectID).Info("Object is not reserved")
        c.JSON(http.StatusBadRequest, gin.H{"error": "Object is not reserved"})
        return
    }

    // Remove the reservation
    object.IsReserved = false
    object.ReservedBy = ""

    // Save the updated object back to DragonflyDB
    data, err := json.Marshal(object)
    if err != nil {
        utils.Log.WithField("objectID", objectID).Error("Failed to marshal updated object")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save unreservation"})
        return
    }

    err = database.RDB.Set(database.Ctx, object.ID, data, 0).Err()
    if err != nil {
        utils.Log.WithField("objectID", objectID).Error("Failed to update object in database")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update object in database"})
        return
    }

    utils.Log.WithField("objectID", object.ID).Info("Object unreserved successfully")
    c.JSON(http.StatusOK, gin.H{"data": object})
}