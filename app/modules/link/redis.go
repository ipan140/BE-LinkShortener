package link

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache interface {
	SetLink(ctx context.Context, code, originalURL string, ttl time.Duration) error
	GetLink(ctx context.Context, code string) (string, error)
	DeleteLink(ctx context.Context, code string) error
	SetLinkMeta(ctx context.Context, code string, link *Link, ttl time.Duration) error
	GetLinkMeta(ctx context.Context, code string) (*Link, error)
}

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) RedisCache {
	return &redisCache{client: client}
}

func (r *redisCache) SetLink(ctx context.Context, code, originalURL string, ttl time.Duration) error {
	key := fmt.Sprintf("link:%s", code)
	return r.client.Set(ctx, key, originalURL, ttl).Err()
}

func (r *redisCache) GetLink(ctx context.Context, code string) (string, error) {
	key := fmt.Sprintf("link:%s", code)
	return r.client.Get(ctx, key).Result()
}

func (r *redisCache) DeleteLink(ctx context.Context, code string) error {
	pipe := r.client.Pipeline()
	pipe.Del(ctx, fmt.Sprintf("link:%s", code))
	pipe.Del(ctx, fmt.Sprintf("link_meta:%s", code))
	_, err := pipe.Exec(ctx)
	return err
}

func (r *redisCache) SetLinkMeta(ctx context.Context, code string, link *Link, ttl time.Duration) error {
	data, err := json.Marshal(link)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("link_meta:%s", code)
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *redisCache) GetLinkMeta(ctx context.Context, code string) (*Link, error) {
	key := fmt.Sprintf("link_meta:%s", code)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var link Link
	if err := json.Unmarshal(data, &link); err != nil {
		return nil, err
	}
	return &link, nil
}
