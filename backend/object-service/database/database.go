package database

import (
	"context"
	"fmt"
	"hexagone/object-service/utils"
	"os"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

func ConnectDatabase() {

	host := os.Getenv("DRAGONFLY_HOST")
	port := os.Getenv("DRAGONFLY_PORT")

	if host == "" || port == "" {
		utils.Log.Fatal("DRAGONFLY_HOST or DRAGONFLY_PORT is not set in the environment variables")
	}

	RDB = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})

	// Test the connection
	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		utils.Log.WithField("error", err.Error()).Error("Failed to connect to DragonflyDB:")
	}
	utils.Log.Info("DragonflyDB connected successfully!")
}
