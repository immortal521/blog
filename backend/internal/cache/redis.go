// Package cache provides redis cache
package cache

import (
	"blog-server/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient interface {
	Raw() *redis.Client
	Close() error
}

type client struct {
	rdb *redis.Client
}

func (c *client) Raw() *redis.Client {
	return c.rdb
}

func (c *client) Close() error {
	return c.rdb.Close()
}

func New(cfg config.RedisConfig) (RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &client{rdb: rdb}, nil
}
