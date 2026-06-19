package handler

import (
	"blog-server/middleware"

	"github.com/labstack/echo/v5"
	"go.uber.org/fx"
)

type Handlers struct {
	fx.In

	Post PostHandler
	Rss  RssHandler
	Auth AuthHandler
	Link LinkHandler
	Model ModelHandler
}

type Middlewares struct {
	fx.In

	Auth *middleware.AuthMiddleware
}

func RegisterRoutes(
	app *echo.Echo,
	h Handlers,
	m Middlewares,
) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	RegisterAuthRoutes(v1, h.Auth)
	RegisterPostRoute(v1, h.Post, m.Auth)
	RegisterRssRoute(v1, h.Rss)
	RegisterLinkRoutes(v1, h.Link)
	RegisterModelRoutes(v1, h.Model)
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
