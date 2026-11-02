package middleware

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

// RequestLogger logs incoming HTTP requests with sub-millisecond latency measurement
func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			latency := time.Since(start)
			req := c.Request()
			res := c.Response()

			log.Printf("[HTTP] %s %s | Status: %d | Latency: %v | IP: %s",
				req.Method,
				req.URL.Path,
				res.Status,
				latency,
				c.RealIP(),
			)

			return err
		}
	}
}
