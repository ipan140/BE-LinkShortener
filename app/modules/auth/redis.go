package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache interface {
	BlacklistToken(ctx context.Context, token string, ttl time.Duration) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
	CacheUserSession(ctx context.Context, userID uint, data string, ttl time.Duration) error
}

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) RedisCache {
	return &redisCache{client: client}
}

func (r *redisCache) BlacklistToken(ctx context.Context, token string, ttl time.Duration) error {
	if r.client == nil {
		return nil
	}
	key := fmt.Sprintf("auth:blacklist:%s", token)
	return r.client.Set(ctx, key, "1", ttl).Err()
}

func (r *redisCache) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	if r.client == nil {
		return false, nil
	}
	key := fmt.Sprintf("auth:blacklist:%s", token)
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *redisCache) CacheUserSession(ctx context.Context, userID uint, data string, ttl time.Duration) error {
	if r.client == nil {
		return nil
	}
	key := fmt.Sprintf("auth:session:%d", userID)
	return r.client.Set(ctx, key, data, ttl).Err()
}
