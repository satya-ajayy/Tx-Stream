package redis

import (
	// Go Internal Packages
	"context"

	// External Packages
	"github.com/redis/go-redis/v9"
)

// Connect connects to the redis db and returns the client.
func Connect(ctx context.Context, uri, password string) (*redis.Client, error) {
	// Configure the Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     uri,      // Redis server address
		Password: password, // Redis password
		DB:       0,        // Default DB
	})

	_, pingErr := rdb.Ping(ctx).Result()
	if pingErr != nil {
		return nil, pingErr
	}
	return rdb, nil
}
