package services_test

import (
	"bytes"
	"encoding/json"
	"hexagone/room-service/src/database"
	"hexagone/room-service/src/models"
	"hexagone/room-service/src/services"
	"hexagone/room-service/src/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

func setupTestServer() {
	// Initialize logger
	utils.InitLogger()
	
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Use in-memory SQLite database
	database.ConnectDatabase(":memory:")
	
	// Set up router with required endpoints
	router = gin.Default()
	router.POST("/rooms", services.CreateRoom)
	router.GET("/rooms", services.ListRooms)
}

func clearDatabase() {
	database.DB.Exec("DELETE FROM rooms")
}

func TestCreateRoom(t *testing.T) {
	setupTestServer()
	defer clearDatabase()

	tests := []struct {
		name         string
		input        map[string]interface{}
		expectedCode int
	}{
		{
			name: "Valid Room Creation",
			input: map[string]interface{}{
				"name":    "Living Room",
				"home_id": 1,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Missing Name",
			input: map[string]interface{}{
				"home_id": 1,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Missing HomeID",
			input: map[string]interface{}{
				"name": "Living Room",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Invalid HomeID Type",
			input: map[string]interface{}{
				"name":    "Living Room",
				"home_id": "invalid",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearDatabase()
			
			jsonInput, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/rooms", bytes.NewBuffer(jsonInput))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedCode, w.Code)
			
			if tt.expectedCode == http.StatusOK {
				var response map[string]models.Room
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response["data"])
				assert.Equal(t, tt.input["name"], response["data"].Name)
				assert.Equal(t, float64(tt.input["home_id"].(int)), float64(response["data"].HomeID))
			}
		})
	}
}

func TestListRooms(t *testing.T) {
	setupTestServer()
	defer clearDatabase()

	// Create test rooms
	testRooms := []services.CreateRoomInput{
		{Name: "Living Room", HomeID: 1},
		{Name: "Kitchen", HomeID: 1},
		{Name: "Bedroom", HomeID: 2},
	}

	for _, room := range testRooms {
		jsonInput, _ := json.Marshal(room)
		req := httptest.NewRequest("POST", "/rooms", bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	tests := []struct {
		name          string
		homeID        string
		expectedCode  int
		expectedCount int
	}{
		{
			name:          "List Rooms for Home 1",
			homeID:        "1",
			expectedCode:  http.StatusOK,
			expectedCount: 2,
		},
		{
			name:          "List Rooms for Home 2",
			homeID:        "2",
			expectedCode:  http.StatusOK,
			expectedCount: 1,
		},
		{
			name:          "Invalid Home ID",
			homeID:        "invalid",
			expectedCode:  http.StatusBadRequest,
			expectedCount: 0,
		},
		{
			name:          "Missing Home ID",
			homeID:        "",
			expectedCode:  http.StatusBadRequest,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/rooms"
			if tt.homeID != "" {
				url += "?home_id=" + tt.homeID
			}
			
			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedCode, w.Code)
			
			if tt.expectedCode == http.StatusOK {
				var response map[string][]models.Room
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response["data"], tt.expectedCount)
				
				if tt.expectedCount > 0 {
					for _, room := range response["data"] {
						homeID := uint(tt.homeID[0] - '0')
						assert.Equal(t, homeID, room.HomeID)
					}
				}
			}
		})
	}
}