package main

import (
	"blog-server/internal/cache"
	"blog-server/internal/config"
	"blog-server/internal/database"
	"blog-server/internal/handler"
	"blog-server/internal/middleware"
	"blog-server/internal/repo"
	"blog-server/internal/router"
	"blog-server/internal/scheduler"
	"blog-server/internal/service"
	"blog-server/pkg/logger"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	redisClient, err := cache.New(cfg.Redis)
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}

	rdb := redisClient.Raw()
	defer redisClient.Close()

	linkRepo := repo.NewLinkRepo()
	linkService := service.NewLinkService(db, linkRepo)
	linkHandler := handler.NewLinkHandler(linkService)
	router.RegisterLinkRoutes(app, linkHandler)

	postRepo := repo.NewPostRepo()
	postService := service.NewPostService(db, rdb, postRepo)
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 监听系统信号
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh // 阻塞直到收到信号
		log.Info("Shutting down gracefully...")
		cancel() // 触发 ctx.Done()
	}()

	scheduler.New(postService).Start(ctx)

	log.Fatal(app.Listen(cfg.Server.GetAddr()).Error())

	<-ctx.Done()
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	_ = app.ShutdownWithContext(shutdownCtx)
}
