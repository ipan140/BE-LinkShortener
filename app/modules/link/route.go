package link

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
	// Public QR endpoint
	g.GET("/links/qr/:code", r.handler.QRCode)

	// Protected endpoints (Supports both Bearer JWT and X-API-Key)
	linkGroup := g.Group("/links", middleware.CombinedAuthMiddleware(r.jwtSecret, r.authService))
	linkGroup.POST("", r.handler.Create)
	linkGroup.GET("", r.handler.List)
	linkGroup.GET("/:id", r.handler.GetByID)
	linkGroup.PUT("/:id", r.handler.Update)
	linkGroup.DELETE("/:id", r.handler.Delete)
}
