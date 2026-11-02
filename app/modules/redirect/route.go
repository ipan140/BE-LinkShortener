package redirect

import (
	"github.com/labstack/echo/v4"
)

type Route struct {
	handler *Handler
}

func NewRoute(handler *Handler) *Route {
	return &Route{
		handler: handler,
	}
}

// RegisterRoutes registers top-level redirect paths for ultra-low latency routing
func (r *Route) RegisterRoutes(e *echo.Echo) {
	// Root level: /:code
	e.GET("/:code", r.handler.Redirect)
	e.HEAD("/:code", r.handler.Redirect)
	e.POST("/:code", r.handler.Redirect)

	// Explicit prefix: /r/:code
	e.GET("/r/:code", r.handler.Redirect)
	e.HEAD("/r/:code", r.handler.Redirect)
	e.POST("/r/:code", r.handler.Redirect)
}
