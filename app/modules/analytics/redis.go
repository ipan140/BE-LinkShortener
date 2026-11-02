package analytics

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"
	"link-shortener/app/modules/redirect"
)

type RedisQueue interface {
	PopBatch(ctx context.Context, batchSize int) ([]redirect.ClickData, error)
}

type redisQueue struct {
	client   *redis.Client
	queueKey string
}

func NewRedisQueue(client *redis.Client, queueKey string) RedisQueue {
	return &redisQueue{
		client:   client,
		queueKey: queueKey,
	}
}

// PopBatch dequeues up to batchSize items from the Redis list
func (r *redisQueue) PopBatch(ctx context.Context, batchSize int) ([]redirect.ClickData, error) {
	if r.client == nil {
		return nil, nil
	}
	items, err := r.client.RPopCount(ctx, r.queueKey, batchSize).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	clicks := make([]redirect.ClickData, 0, len(items))
	for _, raw := range items {
		var click redirect.ClickData
		if err := json.Unmarshal([]byte(raw), &click); err == nil {
			clicks = append(clicks, click)
		}
	}

	return clicks, nil
}
