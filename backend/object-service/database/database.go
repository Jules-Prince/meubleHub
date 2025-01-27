package database

import (
	"context"
	"fmt"

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
		fmt.Println("Failed to connect to DragonflyDB:", err)
	}

	fmt.Println("DragonflyDB connected successfully!")
}
