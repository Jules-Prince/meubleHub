package services

import (
	"hexagone/user-service/src/database"
	"hexagone/user-service/src/models"
	"hexagone/user-service/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	utils.Log.Info("Hashing password")
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to hash password")
	}
	return string(bytes), err
}

// CreateUser handles the creation of a new user
func CreateUser(c *gin.Context) {
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error binding JSON in CreateUser")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.Log.WithFields(logrus.Fields{
		"username": input.Username,
		"email":    input.Email,
	}).Info("Creating user")

	// Hash the password
	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create the user
	user := models.User{Username: input.Username, Email: input.Email, Password: hashedPassword}
	if result := database.DB.Create(&user); result.Error != nil {
		utils.Log.WithFields(logrus.Fields{
			"username": input.Username,
			"email":    input.Email,
			"error":    result.Error.Error(),
		}).Error("Error creating user in the database")
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	utils.Log.WithFields(logrus.Fields{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	}).Info("User created successfully")

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// Login handles user login
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error binding JSON in Login")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.Log.WithFields(logrus.Fields{
		"email": input.Email,
	}).Info("Attempting login")

	// Find the user by email
	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		utils.Log.WithFields(logrus.Fields{
			"email": input.Email,
		}).Warn("Invalid email during login")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Check the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		utils.Log.WithFields(logrus.Fields{
			"email": input.Email,
		}).Warn("Invalid password during login")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	utils.Log.WithFields(logrus.Fields{
		"id":    user.ID,
		"email": user.Email,
	}).Info("Login successful")

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user})
}

// ListUsers retrieves all users
func ListUsers(c *gin.Context) {
	utils.Log.Info("Fetching all users")

	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		utils.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to retrieve users from the database")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	utils.Log.WithFields(logrus.Fields{
		"count": len(users),
	}).Info("Users fetched successfully")

	c.JSON(http.StatusOK, gin.H{"data": users})
}