package router

import (
	"blog-server/internal/handler"
	"blog-server/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterPostRoutes(r fiber.Router, am *middleware.AuthMiddleware, h handler.IPostHandler) {
	group := r.Group("/posts")
	group.Get("/", h.GetPosts)
	group.Get("/meta", h.GetPostIds)
	group.Get("/:id", h.GetPost)

	needAuthGroup := r.Group("/posts", am.Handler())
	needAuthGroup.Post("/test", h.GetPosts)
}
