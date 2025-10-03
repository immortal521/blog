// Package cache provides redis cache
package cache

import (
	"blog-server/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheClient interface {
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

func NewCacheClient(cfg *config.Config) (CacheClient, error) {
	rcfg := cfg.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", rcfg.Host, rcfg.Port),
		Password:     rcfg.Password,
		DB:           rcfg.DB,
		PoolSize:     rcfg.PoolSize,
		MinIdleConns: rcfg.MinIdleConns,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &client{rdb: rdb}, nil
}
