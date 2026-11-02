package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	appErrors "link-shortener/common/errors"
	"link-shortener/common/utils"
)

// RateLimiterConfig defines rate limit settings
type RateLimiterConfig struct {
	Limit  int64         // Maximum allowed requests
	Window time.Duration // Window size
	Prefix string        // Redis key prefix
}

// RedisRateLimiter implements a sliding-window rate limiter powered by Redis
func RedisRateLimiter(rdb *redis.Client, cfg RateLimiterConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if rdb == nil {
				// If Redis is disabled, bypass gracefully
				return next(c)
			}

			ip := utils.GetRealIP(c)
			key := fmt.Sprintf("ratelimit:%s:%s", cfg.Prefix, ip)
			now := time.Now().UnixNano()
			clearBefore := now - cfg.Window.Nanoseconds()

			ctx, cancel := context.WithTimeout(c.Request().Context(), 500*time.Millisecond)
			defer cancel()

			pipe := rdb.Pipeline()
			// Remove entries older than the sliding window
			pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(clearBefore, 10))
			// Add the current request timestamp
			pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})
			// Count entries within window
			countCmd := pipe.ZCard(ctx, key)
			// Set expiration to auto-cleanup the key
			pipe.Expire(ctx, key, cfg.Window)

			_, err := pipe.Exec(ctx)
			if err != nil {
				// On Redis error, fail open to prevent blocking legitimate traffic
				return next(c)
			}

			count := countCmd.Val()
			if count > cfg.Limit {
				c.Response().Header().Set("Retry-After", strconv.FormatInt(int64(cfg.Window.Seconds()), 10))
				return utils.JSONError(c, appErrors.New(http.StatusTooManyRequests, "TOO_MANY_REQUESTS", "rate limit exceeded, please slow down", map[string]any{
					"limit":       cfg.Limit,
					"window_sec":  int64(cfg.Window.Seconds()),
					"retry_after": int64(cfg.Window.Seconds()),
				}))
			}

			return next(c)
		}
	}
}
