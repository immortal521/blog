package router

import (
	"blog-server/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func RegisterLinkRoutes(r *fiber.App, h handler.ILinkHandler) {
	group := r.Group("/api/v1/links")
	group.Get("/", h.GetLinks)
	group.Post("/apply-link", h.ApplyForALinks)
}
