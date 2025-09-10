package router

import (
	"blog-server/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func RegisterPostRoutes(r *fiber.App, h handler.IPostHandler) {
	group := r.Group("/api/v1/posts")
	group.Get("/", h.GetPosts)
	group.Get("/ids", h.GetPostIds)
	group.Get("/:id", h.GetPost)
}
