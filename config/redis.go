package config

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// InitRedis initializes and tests the Redis client with connection pooling
func InitRedis(cfg *RedisConfig) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis server at %s: %w", addr, err)
	}

	log.Printf("[Redis] Successfully connected to Redis pool at %s (PoolSize: %d)", addr, cfg.PoolSize)
	return rdb, nil
}
