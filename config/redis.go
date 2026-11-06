package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// InitRedis initializes and tests the Redis client with connection pooling
func InitRedis(cfg *RedisConfig) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	// Quick ping check with a single connection (no idle pool) to verify availability
	testClient := redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    cfg.Password,
		DB:          cfg.DB,
		DialTimeout: 500 * time.Millisecond,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	if err := testClient.Ping(ctx).Err(); err != nil {
		_ = testClient.Close()
		return nil, fmt.Errorf("redis not reachable at %s: %w", addr, err)
	}
	_ = testClient.Close()

	// Redis is reachable, initialize full pool
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

	log.Printf("[Redis] Successfully connected to Redis pool at %s (PoolSize: %d)", addr, cfg.PoolSize)
	return rdb, nil
}
