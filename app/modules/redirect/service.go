package redirect

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	appErrors "link-shortener/common/errors"
)

type Service interface {
	Resolve(ctx context.Context, shortCode string, click *ClickData) (string, error)
}

type service struct {
	repo  Repository
	cache RedisCache
}

func NewService(repo Repository, cache RedisCache) Service {
	return &service{
		repo:  repo,
		cache: cache,
	}
}

func (s *service) Resolve(ctx context.Context, shortCode string, click *ClickData) (string, error) {
	// 1. Fast Path: Check Redis Cache (Ultra-Fast < 1ms)
	cachedURL, err := s.cache.GetRedirectURL(ctx, shortCode)
	if err == nil && cachedURL != "" {
		// Asynchronously enqueue click data to avoid adding I/O latency to the redirect HTTP response
		go func(c ClickData) {
			bgCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			_ = s.cache.EnqueueClick(bgCtx, &c)
		}(*click)

		return cachedURL, nil
	}

	// 2. Slow Path: Cache Miss - Query Relational Database
	if err != nil && !errors.Is(err, redis.Nil) {
		// Log Redis connection error, continue to DB fallback
	}

	link, err := s.repo.FindActiveLink(ctx, shortCode)
	if err != nil {
		return "", appErrors.Internal("database query error")
	}
	if link == nil {
		return "", appErrors.NotFound("short link not found")
	}

	// Check status and expiration
	if !link.IsActive {
		return "", appErrors.BadRequest("short link is deactivated", nil)
	}
	if link.ExpiresAt != nil && link.ExpiresAt.Before(time.Now()) {
		return "", appErrors.BadRequest("short link has expired", nil)
	}

	// 3. Warm the Redis Cache for future ultra-fast lookups (24 hours TTL)
	_ = s.cache.SetRedirectCache(ctx, shortCode, link.OriginalURL, 24*time.Hour)

	// 4. Asynchronously enqueue click event
	go func(c ClickData) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = s.cache.EnqueueClick(bgCtx, &c)
	}(*click)

	return link.OriginalURL, nil
}
