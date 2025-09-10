package main

import (
	"blog-server/internal/config"
	"blog-server/internal/database"
	"blog-server/internal/handler"
	"blog-server/internal/middleware"
	"blog-server/internal/repo"
	"blog-server/internal/router"
	"blog-server/internal/service"
	"blog-server/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func main() {
	log := logger.Get()

	app := fiber.New(fiber.Config{
		ErrorHandler: handler.ErrorHandler,
	})

	cfg, err := config.Load("./config.yml")
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}

	app.Use(middleware.RequestLogger(log))

	db, err := database.NewDB(cfg.Database.GetDSN())
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}

	linkRepo := repo.NewLinkRepo()
	linkService := service.NewLinkService(db, linkRepo)
	linkHandler := handler.NewLinkHandler(linkService)
	router.RegisterLinkRoutes(app, linkHandler)

	postRepo := repo.NewPostRepo()
	postService := service.NewPostService(db, postRepo)
	postHandler := handler.NewPostHandler(postService)
	router.RegisterPostRoutes(app, postHandler)

	mailService, err := service.NewEmailService()
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}

	authService := service.NewAuthService(mailService)
	authHandler := handler.NewAuthHandler(authService)
	router.RegisterAuthRoutes(app, authHandler)

	log.Fatal(app.Listen(cfg.Server.GetAddr()).Error())
}
