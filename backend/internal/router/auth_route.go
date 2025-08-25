package router

import (
	"blog-server/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(r *fiber.App, h handler.IAuthHandler) {
	group := r.Group("/api/v1/auth")
	group.Post("/captcha", h.SendCaptcha)
}
