// Package router
package router

import (
	"blog-server/internal/handler"

	"github.com/labstack/echo/v4"
)

func RegisterAuthRoutes(r *echo.Group, h handler.IAuthHandler) {
	group := r.Group("/auth")
	group.POST("/captcha", h.SendCaptcha)
	group.POST("/register", h.Register)
	group.POST("/login", h.Login)
	group.POST("/logout", h.Logout)
	group.POST("/refresh", h.Refresh)
}
