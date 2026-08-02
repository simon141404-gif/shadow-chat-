package db

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(url string) (*redis.Client, error) {
	// If no URL provided, return nil (will use in-memory fallback)
	if url == "" || url == "redis://localhost:6379" {
		return nil, nil
	}
	
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return client, nil
}
