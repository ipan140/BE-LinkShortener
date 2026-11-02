package utils

import (
	"net"
	"strings"

	"github.com/labstack/echo/v4"
)

// GetRealIP extracts the client's real public/remote IP address
func GetRealIP(c echo.Context) string {
	req := c.Request()

	// 1. Cloudflare header
	if cfIP := req.Header.Get("CF-Connecting-IP"); cfIP != "" {
		return strings.TrimSpace(cfIP)
	}

	// 2. X-Real-IP
	if realIP := req.Header.Get("X-Real-IP"); realIP != "" {
		return strings.TrimSpace(realIP)
	}

	// 3. X-Forwarded-For (take the first client IP)
	if xff := req.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if ip != "" {
				return ip
			}
		}
	}

	// 4. RemoteAddr
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err == nil {
		return ip
	}

	return req.RemoteAddr
}
