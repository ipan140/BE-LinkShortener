package auth

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
	authGroup := g.Group("/auth")
	authGroup.POST("/register", r.handler.Register)
	authGroup.POST("/login", r.handler.Login)

	// Protected routes
	protected := authGroup.Group("", middleware.CombinedAuthMiddleware(r.jwtSecret, r.handler.service))
	protected.GET("/me", r.handler.Me)
	protected.POST("/api-key/regenerate", r.handler.RegenerateAPIKey)
}
