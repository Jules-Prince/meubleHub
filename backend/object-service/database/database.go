package database

import (
	"context"
	"hexagone/object-service/utils"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

func ConnectDatabase() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // DragonflyDB host and port
		Password: "",               // No password by default
		DB:       0,                // Default DB
	})

	// Test the connection
	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		utils.Log.WithField("error", err.Error()).Error("Failed to connect to DragonflyDB:")
	}
	utils.Log.Info("DragonflyDB connected successfully!")
}
