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
	"blog-server/pkg/validatorx"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func provideConfig() (*config.Config, error) {
	return config.Load("./config.yml")
}

func providerFiberApp(cfg *config.Config, logger *zap.Logger) (*fiber.App, error) {
	var ips []string
	if cfg.App.Environment == "production" {
		cfIPs, err := util.FetchCloudflareIPs()
		if err != nil {
			logger.Error("fetch cloudflare ips failed", zap.Error(err))
			return nil, err
		}
		ips = append(ips, cfIPs...)
	}
	ips = append(ips, "127.0.0.1")

	fiberCfg := fiber.Config{
		EnableTrustedProxyCheck: true,
		ErrorHandler:            handler.ErrorHandler(logger),
		TrustedProxies:          ips,
	}

	app := fiber.New(fiberCfg)
	app.Use(middleware.RequestLogger(logger))

	return app, nil
}

func registerRoutes(app *fiber.App, linkHandler handler.ILinkHandler, postHandler handler.IPostHandler, authHandler handler.IAuthHandler) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	router.RegisterLinkRoutes(v1, linkHandler)
	router.RegisterPostRoutes(v1, postHandler)
	router.RegisterAuthRoutes(v1, authHandler)
}

func runJobsLifecycle(lc fx.Lifecycle, scheduler *scheduler.Scheduler) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			scheduler.Start(ctx)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

func runServerLifecycle(lc fx.Lifecycle, app *fiber.App, cfg *config.Config, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				logger.Info("Server is starting and listening", zap.String("address", cfg.Server.GetAddr()))
				if err := app.Listen(cfg.Server.GetAddr()); err != nil {
					logger.Fatal("Server startup failed", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down server gracefully...")
			timeout := cfg.Server.GracefulShutdown
			if timeout <= 0 {
				timeout = 5 * time.Second
			}
			shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			if err := app.ShutdownWithContext(shutdownCtx); err != nil {
				logger.Error("Server shutdown failed", zap.Error(err))
			} else {
				log.Info("Server has been shut down successfully.")
			}

			return nil
		},
	})
}

func cleanupResources(ls fx.Lifecycle, db database.DB, rc cache.RedisClient) {
	ls.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			db.Close()
			rc.Close()
			return nil
		},
	})
}

func main() {
	provider := fx.Provide(
		logger.NewLogger,
		provideConfig,
		database.NewDB,
		cache.NewRedisClient,
		validatorx.NewValidator,
		providerFiberApp,

		scheduler.NewScheduler,

		// repos
		repo.NewUserRepo,
		repo.NewPostRepo,
		repo.NewLinkRepo,

		// services
		service.NewAuthService,
		service.NewPostService,
		service.NewJwtService,
		service.NewLinkService,
		service.NewEmailService,

		// handlers
		handler.NewAuthHandler,
		handler.NewPostHandler,
		handler.NewLinkHandler,
	)

	invoke := fx.Invoke(
		registerRoutes,
		runServerLifecycle,
		runJobsLifecycle,
		cleanupResources,
	)

	app := fx.New(provider, invoke)

	app.Run()
}
