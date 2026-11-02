package redirect

import (
	"context"
	"errors"
	"net/url"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	appErrors "link-shortener/common/errors"
)

type Service interface {
	Resolve(ctx context.Context, shortCode string, incomingQuery string, password string, click *ClickData) (string, bool, error)
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

// Resolve handles resolution with UTM query forwarding, password checks, and click limits.
// Returns (destinationURL, requiresPassword, error)
func (s *service) Resolve(ctx context.Context, shortCode string, incomingQuery string, password string, click *ClickData) (string, bool, error) {
	// 1. Fast Path: Check Redis Cache (Ultra-Fast < 1ms for non-protected links)
	cachedURL, err := s.cache.GetRedirectURL(ctx, shortCode)
	if err == nil && cachedURL != "" {
		finalURL := mergeQueryParams(cachedURL, incomingQuery)

		// Asynchronously enqueue click event to Redis
		go func(c ClickData) {
			bgCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			_ = s.cache.EnqueueClick(bgCtx, &c)
		}(*click)

		return finalURL, false, nil
	}

	// 2. Slow Path: Cache Miss / Protected Link - Query Database
	if err != nil && !errors.Is(err, redis.Nil) {
		// Log Redis connection error, continue to DB fallback
	}

	link, err := s.repo.FindActiveLink(ctx, shortCode)
	if err != nil {
		return "", false, appErrors.Internal("database query error")
	}
	if link == nil {
		return "", false, appErrors.NotFound("short link not found")
	}

	// Check status
	if !link.IsActive {
		return "", false, appErrors.BadRequest("short link is deactivated", nil)
	}

	// Check expiration date
	if link.ExpiresAt != nil && link.ExpiresAt.Before(time.Now()) {
		return "", false, appErrors.BadRequest("short link has expired", nil)
	}

	// Check maximum click limit
	if link.MaxClicks != nil && link.ClickCount >= *link.MaxClicks {
		return "", false, appErrors.BadRequest("short link has reached its maximum click limit", nil)
	}

	// Check password protection
	if link.IsProtected {
		if password == "" {
			// Signal that password is required
			return "", true, nil
		}

		if err := bcrypt.CompareHashAndPassword([]byte(link.PasswordHash), []byte(password)); err != nil {
			return "", false, appErrors.Unauthorized("incorrect password for this protected link")
		}
	}

	// 3. Warm the Redis Cache if not protected and not near limit
	if !link.IsProtected && link.MaxClicks == nil {
		_ = s.cache.SetRedirectCache(ctx, shortCode, link.OriginalURL, 24*time.Hour)
	}

	// 4. Asynchronously enqueue click event
	go func(c ClickData) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = s.cache.EnqueueClick(bgCtx, &c)
	}(*click)

	finalURL := mergeQueryParams(link.OriginalURL, incomingQuery)
	return finalURL, false, nil
}

// mergeQueryParams forwards UTM tags and other query parameters to the destination URL
func mergeQueryParams(destination string, incomingQuery string) string {
	if incomingQuery == "" {
		return destination
	}

	destURL, err := url.Parse(destination)
	if err != nil {
		return destination
	}

	incomingValues, err := url.ParseQuery(incomingQuery)
	if err != nil {
		return destination
	}

	destValues := destURL.Query()
	for k, vals := range incomingValues {
		for _, v := range vals {
			destValues.Add(k, v)
		}
	}

	destURL.RawQuery = destValues.Encode()
	return destURL.String()
}
