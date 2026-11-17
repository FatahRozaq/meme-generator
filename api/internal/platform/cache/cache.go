// Package cache menyediakan koneksi Redis.
package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/nyxsr/meme-portal/api/internal/config"
)

func NewRedis(ctx context.Context, cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{Addr: cfg.Addr})
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("cache: ping gagal: %w", err)
	}
	return client, nil
}
