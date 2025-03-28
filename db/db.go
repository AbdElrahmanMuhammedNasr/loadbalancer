package db

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

var RedisClient *redis.Client
var ctx = context.Background()

func IntiRedisDB() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:5555", // Change if Redis is running elsewhere
		Password: "",               // No password by default
		DB:       0,                // Default DB
	})

	// Test connection
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis!")

}
