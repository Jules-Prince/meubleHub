package middleware

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func SetupCORS() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost", "http://localhost:80", "http://frontend", "http://frontend:80"},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * 60 * 60, // 12 hours
    })
}