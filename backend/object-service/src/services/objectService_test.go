package services_test

import (
	"bytes"
	"encoding/json"
	"hexagone/object-service/src/database"
	"hexagone/object-service/src/models"
	"hexagone/object-service/src/services"
	"hexagone/object-service/src/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/alicebob/miniredis/v2"
)

var router *gin.Engine
var mr *miniredis.Miniredis

func setupTestServer() error {
	utils.InitLogger()
	gin.SetMode(gin.TestMode)
	
	var err error
	mr, err = miniredis.Run()
	if err != nil {
		return err
	}
	
	os.Setenv("DRAGONFLY_HOST", mr.Host())
	os.Setenv("DRAGONFLY_PORT", mr.Port())
	
	if err := database.ConnectDatabase(); err != nil {
		return err
	}
	
	router = gin.Default()
	router.POST("/objects", services.CreateObject)
	router.GET("/objects", services.ListObjects)
	router.GET("/objects/room", services.ListObjectsByRoom)
	router.PATCH("/objects/:id/reserve", services.ReserveObject)
	router.PATCH("/objects/:id/unreserve", services.UnreserveObject)
	router.GET("/objects/reserved", services.ListReservedObjects)
	
	return nil
}

func cleanupTest() {
	mr.Close()
}

func TestCreateObject(t *testing.T) {
	if err := setupTestServer(); err != nil {
		t.Fatalf("Failed to setup test server: %v", err)
	}
	defer cleanupTest()

	tests := []struct {
		name         string
		input        map[string]interface{}
		expectedCode int
	}{
		{
			name: "Valid Object Creation",
			input: map[string]interface{}{
				"name": "Test Object",
				"type": "furniture",
				"room_id": "room123",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Missing Name",
			input: map[string]interface{}{
				"type": "furniture",
				"room_id": "room123",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Missing Type",
			input: map[string]interface{}{
				"name": "Test Object",
				"room_id": "room123",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Missing Room ID",
			input: map[string]interface{}{
				"name": "Test Object",
				"type": "furniture",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonInput, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/objects", bytes.NewBuffer(jsonInput))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedCode, w.Code)
			
			if tt.expectedCode == http.StatusOK {
				var response map[string]models.Object
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response["data"])
				assert.Equal(t, tt.input["name"], response["data"].Name)
				assert.Equal(t, tt.input["type"], response["data"].Type)
				assert.Equal(t, tt.input["room_id"], response["data"].RoomID)
				assert.False(t, response["data"].IsReserved)
			}
		})
	}
}

func TestListObjectsByRoom(t *testing.T) {
	if err := setupTestServer(); err != nil {
		t.Fatalf("Failed to setup test server: %v", err)
	}
	defer cleanupTest()

	// Create test objects in different rooms
	testObjects := []map[string]interface{}{
		{"name": "Object 1", "type": "furniture", "room_id": "room1"},
		{"name": "Object 2", "type": "electronics", "room_id": "room1"},
		{"name": "Object 3", "type": "furniture", "room_id": "room2"},
	}

	for _, obj := range testObjects {
		jsonInput, _ := json.Marshal(obj)
		req := httptest.NewRequest("POST", "/objects", bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	tests := []struct {
		name           string
		roomID         string
		expectedCode   int
		expectedCount  int
	}{
		{
			name:           "Valid Room ID",
			roomID:         "room1",
			expectedCode:   http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "Empty Room",
			roomID:         "room3",
			expectedCode:   http.StatusOK,
			expectedCount:  0,
		},
		{
			name:           "Missing Room ID",
			roomID:         "",
			expectedCode:   http.StatusBadRequest,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/objects/room"
			if tt.roomID != "" {
				url += "?room_id=" + tt.roomID
			}
			
			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedCode, w.Code)
			
			if tt.expectedCode == http.StatusOK {
				var response map[string][]models.Object
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response["data"], tt.expectedCount)
				
				if tt.expectedCount > 0 {
					for _, obj := range response["data"] {
						assert.Equal(t, tt.roomID, obj.RoomID)
					}
				}
			}
		})
	}
}

func TestUnreserveObject(t *testing.T) {
	if err := setupTestServer(); err != nil {
		t.Fatalf("Failed to setup test server: %v", err)
	}
	defer cleanupTest()

	// Create and reserve an object
	createInput := map[string]interface{}{
		"name": "Test Object",
		"type": "furniture",
		"room_id": "room123",
	}
	jsonInput, _ := json.Marshal(createInput)
	req := httptest.NewRequest("POST", "/objects", bytes.NewBuffer(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	var createResponse map[string]models.Object
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	objectID := createResponse["data"].ID

	// Reserve the object
	reserveInput := map[string]interface{}{
		"userId": "user123",
	}
	jsonInput, _ = json.Marshal(reserveInput)
	req = httptest.NewRequest("PATCH", "/objects/"+objectID+"/reserve", bytes.NewBuffer(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	tests := []struct {
		name         string
		objectID     string
		expectedCode int
	}{
		{
			name:         "Valid Unreservation",
			objectID:     objectID,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Already Unreserved",
			objectID:     objectID,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Non-existent Object",
			objectID:     "nonexistent",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("PATCH", "/objects/"+tt.objectID+"/unreserve", nil)
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedCode, w.Code)
			
			if tt.expectedCode == http.StatusOK {
				var response map[string]models.Object
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.False(t, response["data"].IsReserved)
				assert.Empty(t, response["data"].ReservedBy)
			}
		})
	}
}