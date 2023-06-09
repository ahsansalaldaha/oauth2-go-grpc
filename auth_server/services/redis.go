package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisService - Redis Client
type RedisService struct {
	Client *redis.Client
	ctx    context.Context
}

// Get - Get the redis key if exists
func (rs *RedisService) Get(key string) *redis.StringCmd {
	return rs.Client.Get(rs.ctx, key)
}

// Set - Set the redis key value
func (rs *RedisService) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return rs.Client.Set(rs.ctx, key, value, expiration)
}

// Remember - implements the redis callback functionality
func (rs *RedisService) Remember(key string, duration time.Duration, getDataFunc func() interface{}) interface{} {
	// Check if the data exists in Redis
	val, err := rs.Get(key).Result()
	if err == redis.Nil {
		// Data not found in Redis, call the provided function to get the data
		data := getDataFunc()

		if data != nil {
			// Store the data in Redis with the specified key and expiration duration
			rs.Set(key, data, duration)
		}

		return data
	} else if err != nil {
		// Error occurred while fetching data from Redis
		fmt.Println("Error:", err)
	}

	// Data found in Redis, return it
	return val
}

// NewRedisService - New Redis Client Creator
func NewRedisService(ctx context.Context) *RedisService {
	client := redis.NewClient(&redis.Options{
		Addr:     "cache:6379", // Update with your Redis server address and port
		Password: "",           // Set the password if needed
		DB:       0,            // Select the Redis database
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return &RedisService{
		Client: client,
		ctx:    ctx,
	}
}
