package router

import (
	"blog-server/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterLinkRoutes(r fiber.Router, h handler.ILinkHandler) {
	group := r.Group("/links")
	group.Get("/", h.GetLinks)
	group.Post("/apply-link", h.ApplyForALinks)
	group.Get("/overview", h.GetLinksOverview)
}
