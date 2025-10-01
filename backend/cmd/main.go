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
	"blog-server/pkg/util"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	log := logger.Get()

	cfg, err := config.Load("./config.yml")
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}

	var ips []string
	if cfg.App.Environment == "production" {
		ips, err = util.FetchCloudflareIPs()
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	ips = append(ips, "127.0.0.1")

	fiberCfg := fiber.Config{
		EnableTrustedProxyCheck: true,
		ErrorHandler:            handler.ErrorHandler,
		TrustedProxies:          ips,
	}
	fmt.Println(ips)

	app := fiber.New(fiberCfg)

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

	postRepo := repo.NewPostRepo()
	postService := service.NewPostService(db, rdb, postRepo)
	postHandler := handler.NewPostHandler(postService)

	mailService, err := service.NewEmailService()
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}

	userRepo := repo.NewUserRepo()

	jwtService := service.NewJwtService()
	authService := service.NewAuthService(db, rdb, userRepo, jwtService, mailService)
	authHandler := handler.NewAuthHandler(authService)

	api := app.Group("/api")
	v1 := api.Group("/v1")
	router.RegisterLinkRoutes(v1, linkHandler)
	router.RegisterPostRoutes(v1, postHandler)
	router.RegisterAuthRoutes(v1, authHandler)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	scheduler.New(postService).Start(ctx)

	go func() {
		log.Info(fmt.Sprintf("Server is starting and listening on %s", cfg.Server.GetAddr()))
		if err := app.Listen(cfg.Server.GetAddr()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(fmt.Sprintf("Server listen error: %v", err))
		}
	}()

	<-ctx.Done()
	log.Info("Shutting down server gracefully...")

	// 创建一个带超时的 context 用于关停服务器，防止关停过程无限期阻塞
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Error(fmt.Sprintf("Server shutdown failed: %v", err))
	}

	log.Info("Server has been shut down successfully.")
}
