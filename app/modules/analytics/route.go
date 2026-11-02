package analytics

import (
	"github.com/labstack/echo/v4"
	"link-shortener/common/middleware"
)

type Route struct {
	handler   *Handler
	jwtSecret string
}

func NewRoute(handler *Handler, jwtSecret string) *Route {
	return &Route{
		handler:   handler,
		jwtSecret: jwtSecret,
	}
}

func (r *Route) RegisterRoutes(g *echo.Group) {
	analyticsGroup := g.Group("/analytics", middleware.JWTMiddleware(r.jwtSecret))
	analyticsGroup.GET("/:code", r.handler.GetAnalytics)
}
