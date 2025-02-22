package services

import (
	"hexagone/user-service/src/database"
	"hexagone/user-service/src/models"
	"hexagone/user-service/src/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	AdminKey string `json:"adminKey"`
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

	// Check for admin creation
	isAdmin := false
	adminKey := os.Getenv("ADMIN_KEY")
	if adminKey != "" && input.AdminKey == adminKey {
		isAdmin = true
		utils.Log.Info("Creating admin user")
	}

	// Hash the password
	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create the user
	user := models.User{
		Username: input.Username, 
		Email: input.Email, 
		Password: hashedPassword,
		IsAdmin: isAdmin,
	}
	
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
		"isAdmin":  user.IsAdmin,
	}).Info("User created successfully")

	// Remove password from response
	user.Password = ""
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

	// Create a safe user copy without password
	safeUser := models.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		IsAdmin:  user.IsAdmin,
	}

	utils.Log.WithFields(logrus.Fields{
		"id":      user.ID,
		"email":   user.Email,
		"isAdmin": user.IsAdmin,
	}).Info("Login successful")

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": safeUser})
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

func GetUser(c *gin.Context) {
    userID := c.Param("id")
    
    utils.Log.WithFields(logrus.Fields{
        "userID": userID,
    }).Info("Fetching user by ID")

    var user models.User
    if err := database.DB.First(&user, userID).Error; err != nil {
        utils.Log.WithFields(logrus.Fields{
            "userID": userID,
            "error":  err.Error(),
        }).Error("Failed to retrieve user from the database")
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    utils.Log.WithFields(logrus.Fields{
        "userID":   user.ID,
        "username": user.Username,
    }).Info("User fetched successfully")

    c.JSON(http.StatusOK, gin.H{"data": user})
}