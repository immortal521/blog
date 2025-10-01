// Package router
package router

import (
	"blog-server/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(r fiber.Router, h handler.IAuthHandler) {
	group := r.Group("/auth")
	group.Post("/captcha", h.SendCaptcha)
	group.Post("/register", h.Register)
	group.Post("/login", h.Login)
	group.Post("/logout", h.Logout)
}
