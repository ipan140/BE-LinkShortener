package redirect

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	appErrors "link-shortener/common/errors"
)

type Service interface {
	Resolve(ctx context.Context, shortCode string, incomingQuery string, password string, click *ClickData) (string, bool, error)
	GetLinkDetails(ctx context.Context, shortCode string) (*RedirectLink, error)
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

// Resolve handles resolution with Smart Device Routing, UTM forwarding, password checks, and click limits.
// Returns (destinationURL, requiresPassword, error)
func (s *service) Resolve(ctx context.Context, shortCode string, incomingQuery string, password string, click *ClickData) (string, bool, error) {
	// 1. Fast Path: Check Redis Cache (Ultra-Fast < 1ms for non-protected, non-device-routed links)
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

	// 2. Slow Path: Cache Miss / Protected Link / Dynamic Device Routing - Query Database
	if err != nil && !errors.Is(err, redis.Nil) {
		// Redis error, proceed to DB fallback
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
			return "", true, nil
		}

		if err := bcrypt.CompareHashAndPassword([]byte(link.PasswordHash), []byte(password)); err != nil {
			return "", false, appErrors.Unauthorized("incorrect password for this protected link")
		}
	}

	// 3. Smart Device Routing (iOS / Android / Default)
	chosenDestination := link.OriginalURL
	osLower := strings.ToLower(click.OS)
	if (strings.Contains(osLower, "ios") || strings.Contains(osLower, "ipad")) && link.IOSURL != "" {
		chosenDestination = link.IOSURL
	} else if strings.Contains(osLower, "android") && link.AndroidURL != "" {
		chosenDestination = link.AndroidURL
	}

	// 4. Warm the Redis Cache ONLY if the link is purely static and unconstrained
	canCache := !link.IsProtected && link.MaxClicks == nil && link.IOSURL == "" && link.AndroidURL == ""
	if canCache {
		_ = s.cache.SetRedirectCache(ctx, shortCode, link.OriginalURL, 24*time.Hour)
	}

	// 5. Asynchronously enqueue click event
	go func(c ClickData) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = s.cache.EnqueueClick(bgCtx, &c)
	}(*click)

	finalURL := mergeQueryParams(chosenDestination, incomingQuery)
	return finalURL, false, nil
}

func (s *service) GetLinkDetails(ctx context.Context, shortCode string) (*RedirectLink, error) {
	return s.repo.FindActiveLink(ctx, shortCode)
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
