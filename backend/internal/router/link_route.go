package router

import (
	"blog-server/internal/handler"

	"github.com/labstack/echo/v4"
)

func RegisterLinkRoutes(r *echo.Group, h handler.ILinkHandler) {
	group := r.Group("/links")
	group.GET("", h.GetLinks)
	group.POST("/apply-link", h.ApplyForALinks)
}
