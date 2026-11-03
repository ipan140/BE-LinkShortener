package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTPRequestsTotal tracks the total number of HTTP requests
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed by method, path, and status code.",
		},
		[]string{"method", "path", "status"},
	)

	// HTTPRequestDuration tracks latency in seconds with sub-millisecond buckets for <10ms verification
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Histogram of request latencies in seconds.",
			Buckets: []float64{
				0.0005, // 0.5ms
				0.0010, // 1ms
				0.0020, // 2ms
				0.0050, // 5ms
				0.0100, // 10ms (Target SLA threshold)
				0.0250, // 25ms
				0.0500, // 50ms
				0.1000, // 100ms
				0.2500, // 250ms
				0.5000, // 500ms
				1.0000, // 1s
			},
		},
		[]string{"method", "path", "status"},
	)

	// InFlightRequests tracks currently active concurrent requests
	InFlightRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_in_flight_requests",
			Help: "Current number of in-flight concurrent HTTP requests.",
		},
	)

	// RedirectsTotal tracks total redirect events categorized by status
	RedirectsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shortener_redirects_total",
			Help: "Total shortlink redirects processed by outcome (hit_cache, hit_db, expired, limit_reached, not_found).",
		},
		[]string{"outcome"},
	)
)

// PrometheusMetrics returns an Echo middleware that records Prometheus telemetry for every request
func PrometheusMetrics() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Path()
			if path == "/metrics" || path == "/health" {
				return next(c)
			}

			InFlightRequests.Inc()
			defer InFlightRequests.Dec()

			start := time.Now()
			err := next(c)
			if err != nil {
				c.Error(err)
			}

			duration := time.Since(start).Seconds()
			status := strconv.Itoa(c.Response().Status)
			method := c.Request().Method

			// Sanitize path for Prometheus label cardinality
			if path == "" {
				path = c.Request().URL.Path
			}

			HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
			HTTPRequestDuration.WithLabelValues(method, path, status).Observe(duration)

			return err
		}
	}
}
