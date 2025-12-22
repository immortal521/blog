// Package router
package router

import (
	"blog-server/internal/handler"
	"blog-server/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App,
	linkHandler handler.ILinkHandler,
	postHandler handler.IPostHandler,
	authHandler handler.IAuthHandler,
	rssHandler handler.IRssHandler,
	modelHandler handler.IModelHandler,
	imageHandler handler.IImageHandler,
	am *middleware.AuthMiddleware,
) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	RegisterLinkRoutes(v1, linkHandler)
	RegisterPostRoutes(v1, am, postHandler)
	RegisterAuthRoutes(v1, authHandler)
	RegisterRssRoutes(v1, rssHandler)
	RegisterModelRoutes(v1, modelHandler)
	RegisterFileRoutes(v1, imageHandler)
}
