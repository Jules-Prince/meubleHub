package middleware

import (
	"hexagone/user-service/src/models"
	"hexagone/user-service/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get user from context (after auth middleware)
        userInterface, exists := c.Get("user")
        if !exists {
            utils.Log.Warn("No user found in context")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
            c.Abort()
            return
        }

        // Type assert the user interface to your User model
        user, ok := userInterface.(models.User)
        if !ok {
            utils.Log.Error("Failed to convert user interface to User model")
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
            c.Abort()
            return
        }

        // Check if user is admin
        if !user.IsAdmin {
            utils.Log.Warn("Non-admin user attempted admin action")
            c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
            c.Abort()
            return
        }

        c.Next()
    }
}