package link

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
	linkGroup := g.Group("/links", middleware.JWTMiddleware(r.jwtSecret))
	linkGroup.POST("", r.handler.Create)
	linkGroup.GET("", r.handler.List)
	linkGroup.GET("/:id", r.handler.GetByID)
	linkGroup.PUT("/:id", r.handler.Update)
	linkGroup.DELETE("/:id", r.handler.Delete)
}
