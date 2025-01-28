package services_test

import (
	"bytes"
	"encoding/json"
	"hexagone/user-service/src/database"
	"hexagone/user-service/src/models"
	"hexagone/user-service/src/services"
	"hexagone/user-service/src/utils"
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
	router.POST("/users", services.CreateUser)
	router.POST("/login", services.Login)
	router.GET("/users", services.ListUsers)
}

func clearDatabase() {
	database.DB.Exec("DELETE FROM users")
}

func TestCreateUser(t *testing.T) {
	setupTestServer()
	defer clearDatabase()

	// First, create a user that we'll try to duplicate
	initialUser := map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}
	
	jsonInput, _ := json.Marshal(initialUser)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	tests := []struct {
		name         string
		input        map[string]interface{}
		expectedCode int
	}{
		{
			name: "Valid User Creation",
			input: map[string]interface{}{
				"username": "newuser",
				"email":    "new@example.com",
				"password": "password123",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Missing Username",
			input: map[string]interface{}{
				"email":    "test2@example.com",
				"password": "password123",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Invalid Email",
			input: map[string]interface{}{
				"username": "testuser3",
				"email":    "invalid-email",
				"password": "password123",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Missing Password",
			input: map[string]interface{}{
				"username": "testuser4",
				"email":    "test4@example.com",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Duplicate Email",
			input: map[string]interface{}{
				"username": "testuser5",
				"email":    "test@example.com", // Same email as initialUser
				"password": "password123",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Duplicate Username",
			input: map[string]interface{}{
				"username": "testuser", // Same username as initialUser
				"email":    "different@example.com",
				"password": "password123",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonInput, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonInput))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedCode, w.Code, "Test case: "+tt.name)
			
			if tt.expectedCode == http.StatusOK {
				var response map[string]models.User
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response["data"])
				assert.Equal(t, tt.input["username"], response["data"].Username)
				assert.Equal(t, tt.input["email"], response["data"].Email)
				assert.NotEqual(t, tt.input["password"], response["data"].Password) // Password should be hashed
			}
		})
	}
}

func TestLogin(t *testing.T) {
	setupTestServer()
	defer clearDatabase()

	// Create a test user first
	testUser := map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}
	
	jsonInput, _ := json.Marshal(testUser)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	tests := []struct {
		name         string
		input        map[string]interface{}
		expectedCode int
	}{
		{
			name: "Valid Login",
			input: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Invalid Password",
			input: map[string]interface{}{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "Invalid Email",
			input: map[string]interface{}{
				"email":    "nonexistent@example.com",
				"password": "password123",
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "Missing Email",
			input: map[string]interface{}{
				"password": "password123",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Missing Password",
			input: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonInput, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonInput))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedCode, w.Code)
			
			if tt.expectedCode == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Login successful", response["message"])
				
				user := response["user"].(map[string]interface{})
				assert.Equal(t, testUser["email"], user["email"])
				assert.Equal(t, testUser["username"], user["username"])
			}
		})
	}
}

func TestListUsers(t *testing.T) {
	setupTestServer()
	defer clearDatabase()

	// Create test users
	testUsers := []map[string]interface{}{
		{
			"username": "user1",
			"email":    "user1@example.com",
			"password": "password123",
		},
		{
			"username": "user2",
			"email":    "user2@example.com",
			"password": "password123",
		},
	}

	for _, user := range testUsers {
		jsonInput, _ := json.Marshal(user)
		req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	t.Run("List All Users", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string][]models.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response["data"], len(testUsers))
		
		// Verify user data
		for i, user := range response["data"] {
			assert.Equal(t, testUsers[i]["username"], user.Username)
			assert.Equal(t, testUsers[i]["email"], user.Email)
			assert.NotEqual(t, testUsers[i]["password"], user.Password) // Password should be hashed
		}
	})
}