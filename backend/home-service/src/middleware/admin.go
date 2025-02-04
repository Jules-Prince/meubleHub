package middleware

import (
	"encoding/json"
	"hexagone/home-service/src/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type User struct {
	ID      uint `json:"id"`
	IsAdmin bool `json:"isAdmin"`
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			utils.Log.Warn("No user ID found in header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Use the service name instead of localhost
		userServiceURL := "http://user-service:" + os.Getenv("USER_PORT")

		utils.Log.WithFields(logrus.Fields{
			"userID":         userID,
			"userServiceURL": userServiceURL,
		}).Info("Attempting to verify admin status")

		resp, err := http.Get(userServiceURL + "/users/" + userID)
		if err != nil {
			utils.Log.WithError(err).Error("Failed to contact user service")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify admin status"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			utils.Log.Warn("Failed to get user details")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			c.Abort()
			return
		}

		var response struct {
			Data User `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			utils.Log.WithError(err).Error("Failed to decode user response")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify admin status"})
			c.Abort()
			return
		}

		if !response.Data.IsAdmin {
			utils.Log.Warn("Non-admin user attempted admin action")
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
