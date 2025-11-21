package router

import (
	"blog-server/internal/handler"
	"blog-server/internal/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterPostRoutes(r *echo.Group, am *middleware.AuthMiddleware, h handler.IPostHandler) {
	group := r.Group("/posts")
	group.GET("", h.GetPosts)
	group.GET("/meta", h.GetPostIds)
	group.GET("/:id", h.GetPost)
}
