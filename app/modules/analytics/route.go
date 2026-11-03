package analytics

import (
	"github.com/labstack/echo/v4"
	"link-shortener/common/middleware"
)

type Route struct {
	handler     *Handler
	jwtSecret   string
	authService any
}

func NewRoute(handler *Handler, jwtSecret string, authService any) *Route {
	return &Route{
		handler:     handler,
		jwtSecret:   jwtSecret,
		authService: authService,
	}
}

func (r *Route) RegisterRoutes(g *echo.Group) {
	// Public analytics endpoint (Dashboard integration)
	g.GET("/analytics/public/:code", r.handler.PublicGetAnalytics)

	// Protected endpoints (Supports both Bearer JWT and X-API-Key)
	analyticsGroup := g.Group("/analytics", middleware.CombinedAuthMiddleware(r.jwtSecret, r.authService))
	analyticsGroup.GET("/:code", r.handler.GetAnalytics)
	analyticsGroup.GET("/:code/export.csv", r.handler.ExportCSV)
}
