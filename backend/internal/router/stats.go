package router

import (
	"blog-server/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterStatsRoutes(r fiber.Router, h handler.IStatsHandler) {
	group := r.Group("/stats")
	group.Get("/dashboard", h.GetDashboardStats)
}
