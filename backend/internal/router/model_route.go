package router

import (
	"blog-server/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterModelRoutes(r fiber.Router, h handler.IModelHandler) {
	group := r.Group("/model")
	group.Post("/summarize", h.CreateSummarySession)
	group.Get("/summarize/:sessionId", h.SummaryStream)
}
