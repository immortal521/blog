package handler

import (
	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
)

func RegisterRoutes(app *fiber.App, ph PostHandler, rh RssHandler, ah AuthHandler, lh LinkHandler, mh ModelHandler) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	RegisterPostRoute(v1, ph)
	RegisterRssRoute(v1, rh)
	RegisterAuthRoutes(v1, ah)
	RegisterLinkRoutes(v1, lh)
	RegisterModelRoutes(v1, mh)
}

func Module() fx.Option {
	return fx.Module(
		"handler",
		fx.Provide(
			NewPostHandler,
			NewRssHandler,
			NewAuthHandler,
			NewLinkHandler,
			NewModelHandler,
		),
		fx.Invoke(
			RegisterRoutes,
		),
	)
}
