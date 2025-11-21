package router

import (
	"blog-server/internal/handler"

	"github.com/labstack/echo/v4"
)

func RegisterRssRoutes(r *echo.Group, h handler.IRssHandler) {
	group := r.Group("/rss")
	group.GET("", h.Subscribe)
}
