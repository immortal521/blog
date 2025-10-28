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
	"go.uber.org/fx"
)

func provideConfig() (*config.Config, error) {
	return config.Load("./config.yml")
}

func providerFiberApp(cfg *config.Config, log logger.Logger) (*fiber.App, error) {
	var ips []string
	if cfg.App.Environment == "production" {
		cfIPs, err := util.FetchCloudflareIPs()
		if err != nil {
			log.Error("fetch cloudflare ips failed", logger.Error(err))
			return nil, err
		}
		ips = append(ips, cfIPs...)
	}
	ips = append(ips, "127.0.0.1")

	fiberCfg := fiber.Config{
		EnableTrustedProxyCheck: true,
		ErrorHandler:            handler.ErrorHandler(log),
		TrustedProxies:          ips,
	}

	app := fiber.New(fiberCfg)
	app.Use(middleware.RequestLogger(log))

	return app, nil
}

func registerRoutes(app *fiber.App,
	linkHandler handler.ILinkHandler,
	postHandler handler.IPostHandler,
	authHandler handler.IAuthHandler,
	rssHandler handler.IRssHandler,
	am *middleware.AuthMiddleware) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	router.RegisterLinkRoutes(v1, linkHandler)
	router.RegisterPostRoutes(v1, am, postHandler)
	router.RegisterAuthRoutes(v1, authHandler)
	router.RegisterRssRoutes(v1, rssHandler)
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
		cache.NewCacheClient,
		validatorx.NewValidator,
		providerFiberApp,

		scheduler.NewScheduler,
		middleware.NewAuthMiddleware,

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
		service.NewRssService,

		// handlers
		handler.NewAuthHandler,
		handler.NewPostHandler,
		handler.NewLinkHandler,
		handler.NewRssHandler,
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
