package router

import (
	"blog-server/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterPostRoutes(r fiber.Router, h handler.IPostHandler) {
	group := r.Group("/posts")
	group.Get("/", h.GetPosts)
	group.Get("/meta", h.GetPostIds)
	group.Get("/:id", h.GetPost)
}
