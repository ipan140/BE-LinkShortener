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
	// Root level: /:code (e.g. https://domain.com/abc12)
	e.GET("/:code", r.handler.Redirect)

	// Explicit prefix: /r/:code (e.g. https://domain.com/r/abc12)
	e.GET("/r/:code", r.handler.Redirect)
}
