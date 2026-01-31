package main

import (
	"context"
	"time"

	"blog-server/cache"
	"blog-server/config"
	"blog-server/database"
	"blog-server/handler"
	"blog-server/logger"
	"blog-server/middleware"
	"blog-server/repository"
	"blog-server/router"
	"blog-server/scheduler"
	"blog-server/service"
	"blog-server/storage"
	"blog-server/utils"
	"blog-server/validatorx"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

func provideConfig() (*config.Config, error) {
	_, err := config.Load("./config.yml")
	if err != nil {
		return nil, err
	}
	return config.Get(), nil
}

func providerFiberApp(cfg *config.Config, log logger.Logger) (*fiber.App, error) {
	var ips []string
	if cfg.App.Environment == config.EnvProd {
		cfIPs, err := utils.FetchCloudflareIPs()
		if err != nil {
			log.Error("fetch cloudflare ips failed", logger.Error(err))
			return nil, err
		}
		ips = append(ips, cfIPs...)
	}
	ips = append(ips, "127.0.0.1")

	fiberCfg := fiber.Config{
		EnableTrustedProxyCheck: true,
		ErrorHandler:            handler.ErrorHandler(log, cfg),
		TrustedProxies:          ips,
		ProxyHeader:             fiber.HeaderXForwardedFor,
		BodyLimit:               10 * 1024 * 1024,
	}

	app := fiber.New(fiberCfg)
	app.Use(middleware.RequestLogger(log))

	return app, nil
}

// scheduler start
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

func runServerLifecycle(lc fx.Lifecycle, app *fiber.App, cfg *config.Config, log logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				log.Info("Server is starting")
				if err := app.Listen(cfg.Server.GetAddr()); err != nil {
					log.Fatal("Server startup failed", logger.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Shutting down server gracefully...")
			timeout := cfg.Server.GracefulShutdown
			if timeout <= 0 {
				timeout = 5 * time.Second
			}
			shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			if err := app.ShutdownWithContext(shutdownCtx); err != nil {
				log.Error("Server shutdown failed", logger.Error(err))
			} else {
				log.Info("Server has been shut down successfully.")
			}

			return nil
		},
	})
}

func cleanupResources(ls fx.Lifecycle, db database.DB, rc cache.CacheClient) {
	ls.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			_ = db.Close()
			_ = rc.Close()
			return nil
		},
	})
}

func main() {
	provider := fx.Provide(
		logger.NewLogger,
		provideConfig,
		database.New,
		cache.NewCacheClient,
		validatorx.NewValidator,
		providerFiberApp,

		scheduler.NewScheduler,
		middleware.NewAuthMiddleware,

		// repos
		repository.NewUserRepo,
		repository.NewPostRepo,
		repository.NewLinkRepo,

		// services
		service.NewAuthService,
		service.NewPostService,
		service.NewJwtService,
		service.NewLinkService,
		service.NewEmailService,
		service.NewRssService,
		service.NewModelService,
		service.NewImageService,
		service.NewStatsService,

		// handlers
		handler.NewAuthHandler,
		handler.NewPostHandler,
		handler.NewLinkHandler,
		handler.NewRssHandler,
		handler.NewModelHandler,
		handler.NewImageHandler,
		handler.NewStatsHandler,

		storage.NewS3Storage,
	)

	invoke := fx.Invoke(
		router.RegisterRoutes,
		runServerLifecycle,
		runJobsLifecycle,
		cleanupResources,
	)

	app := fx.New(provider, invoke)

	app.Run()
}
