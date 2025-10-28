package router

import (
	"blog-server/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterRssRoutes(r fiber.Router, h handler.IRssHandler) {
	group := r.Group("/rss")
	group.Get("/", h.Subscribe)
}
