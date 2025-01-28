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
	// Initialize logger
	utils.InitLogger()
	
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Create mock Redis server
	var err error
	mr, err = miniredis.Run()
	if err != nil {
		return err
	}
	
	// Set environment variables for database connection
	os.Setenv("DRAGONFLY_HOST", mr.Host())
	os.Setenv("DRAGONFLY_PORT", mr.Port())
	
	// Initialize database connection
	if err := database.ConnectDatabase(); err != nil {
		return err
	}
	
	// Set up router
	router = gin.Default()
	router.POST("/objects", services.CreateObject)
	router.GET("/objects", services.ListObjects)
	router.PATCH("/objects/:id/reserve", services.ReserveObject)
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
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Missing Name",
			input: map[string]interface{}{
				"type": "furniture",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Missing Type",
			input: map[string]interface{}{
				"name": "Test Object",
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
				assert.False(t, response["data"].IsReserved)
			}
		})
	}
}

func TestReserveObject(t *testing.T) {
	if err := setupTestServer(); err != nil {
		t.Fatalf("Failed to setup test server: %v", err)
	}
	defer cleanupTest()

	// First create an object
	createInput := map[string]interface{}{
		"name": "Test Object",
		"type": "furniture",
	}
	jsonInput, _ := json.Marshal(createInput)
	req := httptest.NewRequest("POST", "/objects", bytes.NewBuffer(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	var createResponse map[string]models.Object
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	objectID := createResponse["data"].ID

	tests := []struct {
		name         string
		objectID     string
		userID       string
		expectedCode int
	}{
		{
			name:         "Valid Reservation",
			objectID:     objectID,
			userID:       "user123",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Already Reserved",
			objectID:     objectID,
			userID:       "user456",
			expectedCode: http.StatusConflict,
		},
		{
			name:         "Non-existent Object",
			objectID:     "nonexistent",
			userID:       "user123",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reserveInput := map[string]interface{}{
				"userId": tt.userID,
			}
			jsonInput, _ := json.Marshal(reserveInput)
			req := httptest.NewRequest("PATCH", "/objects/"+tt.objectID+"/reserve", bytes.NewBuffer(jsonInput))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedCode, w.Code)
			
			if tt.expectedCode == http.StatusOK {
				var response map[string]models.Object
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.True(t, response["data"].IsReserved)
				assert.Equal(t, tt.userID, response["data"].ReservedBy)
			}
		})
	}
}

func TestListObjects(t *testing.T) {
	if err := setupTestServer(); err != nil {
		t.Fatalf("Failed to setup test server: %v", err)
	}
	defer cleanupTest()

	// Create test objects
	testObjects := []map[string]interface{}{
		{"name": "Object 1", "type": "furniture"},
		{"name": "Object 2", "type": "electronics"},
	}

	for _, obj := range testObjects {
		jsonInput, _ := json.Marshal(obj)
		req := httptest.NewRequest("POST", "/objects", bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	t.Run("List All Objects", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/objects", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string][]models.Object
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response["data"], len(testObjects))
	})
}

func TestListReservedObjects(t *testing.T) {
	if err := setupTestServer(); err != nil {
		t.Fatalf("Failed to setup test server: %v", err)
	}
	defer cleanupTest()

	// Create and reserve test objects
	testObjects := []map[string]interface{}{
		{"name": "Reserved Object 1", "type": "furniture"},
		{"name": "Reserved Object 2", "type": "electronics"},
		{"name": "Unreserved Object", "type": "furniture"},
	}

	for i, obj := range testObjects {
		// Create object
		jsonInput, _ := json.Marshal(obj)
		req := httptest.NewRequest("POST", "/objects", bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		var createResponse map[string]models.Object
		json.Unmarshal(w.Body.Bytes(), &createResponse)
		
		// Reserve first two objects
		if i < 2 {
			reserveInput := map[string]interface{}{
				"userId": "user123",
			}
			jsonInput, _ := json.Marshal(reserveInput)
			req = httptest.NewRequest("PATCH", "/objects/"+createResponse["data"].ID+"/reserve", bytes.NewBuffer(jsonInput))
			req.Header.Set("Content-Type", "application/json")
			w = httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	}

	t.Run("List Reserved Objects", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/objects/reserved", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string][]models.Object
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response["data"], 2)
		
		for _, obj := range response["data"] {
			assert.True(t, obj.IsReserved)
			assert.Equal(t, "user123", obj.ReservedBy)
		}
	})
}