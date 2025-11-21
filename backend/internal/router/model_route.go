package router

import (
	"blog-server/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterModelRoutes(r fiber.Router, h handler.IModelHandler) {
	group := r.Group("/model")
	group.Post("/summary", h.ProcessContent)
	group.Get("/sse", h.SSE)
}
