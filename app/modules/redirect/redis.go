package redirect

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache interface {
	GetRedirectURL(ctx context.Context, shortCode string) (string, error)
	SetRedirectCache(ctx context.Context, shortCode, originalURL string, ttl time.Duration) error
	EnqueueClick(ctx context.Context, click *ClickData) error
}

type redisCache struct {
	client   *redis.Client
	queueKey string
}

func NewRedisCache(client *redis.Client, queueKey string) RedisCache {
	return &redisCache{
		client:   client,
		queueKey: queueKey,
	}
}

// GetRedirectURL executes an ultra-fast Redis GET (typical latency < 1 ms)
func (r *redisCache) GetRedirectURL(ctx context.Context, shortCode string) (string, error) {
	key := fmt.Sprintf("link:%s", shortCode)
	return r.client.Get(ctx, key).Result()
}

// SetRedirectCache writes the original URL into Redis with TTL
func (r *redisCache) SetRedirectCache(ctx context.Context, shortCode, originalURL string, ttl time.Duration) error {
	key := fmt.Sprintf("link:%s", shortCode)
	return r.client.Set(ctx, key, originalURL, ttl).Err()
}

// EnqueueClick pushes serialized click event into Redis List queue for asynchronous background processing
func (r *redisCache) EnqueueClick(ctx context.Context, click *ClickData) error {
	payload, err := json.Marshal(click)
	if err != nil {
		return err
	}

	// Use LPush for atomic fast queue insertion
	err = r.client.LPush(ctx, r.queueKey, payload).Err()
	if err != nil {
		log.Printf("[Redirect Redis] Warning: failed to enqueue click: %v", err)
	}
	return err
}
