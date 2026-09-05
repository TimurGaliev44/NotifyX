package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
}

func New(ctx context.Context, conn string) (*Client, error) {

	opts, err := redis.ParseURL(conn)
	if err != nil {
		return nil, fmt.Errorf("failed parse redis url: %w", err)
	}

	redisClient := redis.NewClient(opts)

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed ping redis client: %w", err)
	}

	return &Client{redisClient}, nil
}

func (r *Client) Close() {
	r.client.Close()
}
