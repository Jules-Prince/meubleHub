package services_test

import (
	"bytes"
	"encoding/json"
	"hexagone/home-service/src/database"
	"hexagone/home-service/src/models"
	"hexagone/home-service/src/utils"
	"hexagone/home-service/src/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

// TestHomeInput mirrors the CreateHomeInput struct from services package
type TestHomeInput struct {
	Name string `json:"name"`
}

func setupTest() {
	// Initialize logger
	utils.InitLogger()
	
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Connect to in-memory database
	database.ConnectDatabase(":memory:")
	
	// Set up router
	router = gin.Default()
	router.POST("/homes", services.CreateHome)
	router.GET("/homes", services.ListHomes)
}

func clearDatabase() {
	database.DB.Exec("DELETE FROM homes")
}

func TestCreateHome(t *testing.T) {
	setupTest()
	defer clearDatabase()

	tests := []struct {
		name           string
		input         TestHomeInput
		expectedCode   int
		expectedError  bool
	}{
		{
			name:          "Valid Home Creation",
			input:         TestHomeInput{Name: "Test House"},
			expectedCode:  http.StatusOK,
			expectedError: false,
		},
		{
			name:          "Empty Name",
			input:         TestHomeInput{Name: ""},
			expectedCode:  http.StatusBadRequest,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear database before each sub-test
			clearDatabase()
			
			jsonInput, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/homes", bytes.NewBuffer(jsonInput))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedCode, w.Code)
			
			if !tt.expectedError {
				var response map[string]models.Home
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response["data"])
				assert.Equal(t, tt.input.Name, response["data"].Name)
			}
		})
	}
}

func TestListHomes(t *testing.T) {
	setupTest()
	defer clearDatabase()

	// Clear database at start
	clearDatabase()

	// Create test homes
	testHomes := []TestHomeInput{
		{Name: "House 1"},
		{Name: "House 2"},
	}
	
	// Add test homes
	for _, home := range testHomes {
		jsonInput, _ := json.Marshal(home)
		req := httptest.NewRequest("POST", "/homes", bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Failed to create test home")
	}
	
	t.Run("List All Homes", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/homes", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string][]models.Home
		err := json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.NoError(t, err)
		assert.NotEmpty(t, response["data"])
		assert.Len(t, response["data"], len(testHomes), "Should have exactly the number of homes we created")
		
		// Verify the homes are the ones we created
		homeNames := make(map[string]bool)
		for _, home := range response["data"] {
			homeNames[home.Name] = true
		}
		
		for _, expectedHome := range testHomes {
			assert.True(t, homeNames[expectedHome.Name], "Should find home with name: "+expectedHome.Name)
		}
	})
}