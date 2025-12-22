package router

import (
	"blog-server/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterFileRoutes(r fiber.Router, ih handler.IImageHandler) {
	uploadGroup := r.Group("/upload")
	uploadGroup.Post("/image", ih.Upload)

	downloadGroup := r.Group("/download")
	downloadGroup.Get("/image/:key", ih.Download)

	deleteGroup := r.Group("/delete")
	deleteGroup.Delete("/image/:key", ih.Delete)
}
