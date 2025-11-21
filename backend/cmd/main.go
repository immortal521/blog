package main

import (
	"context"
	"time"

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
	"blog-server/pkg/validatorx"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func provideConfig() (*config.Config, error) {
	_, err := config.Load("./config.yml")
	if err != nil {
		return nil, err
	}
	return config.Get(), nil
}

func providerEchoApp(cfg *config.Config, log logger.Logger) (*echo.Echo, error) {
	// var ips []string
	// if cfg.App.Environment == config.EnvProd {
	// 	cfIPs, err := utils.FetchCloudflareIPs()
	// 	if err != nil {
	// 		log.Error("fetch cloudflare ips failed", logger.Error(err))
	// 		return nil, err
	// 	}
	// 	ips = append(ips, cfIPs...)
	// }
	// ips = append(ips, "127.0.0.1")

	e := echo.New()

	// 设置全局错误处理器
	e.HTTPErrorHandler = handler.ErrorHandler(log, cfg)

	// 设置请求日志中间件
	e.Use(middleware.RequestLogger(log))

	// // 信任代理（Echo 通过 RealIP 支持 X-Forwarded-For）
	// e.IPExtractor = echo.ExtractIPFromXFFHeader(ips...)

	return e, nil
}

func registerRoutes(
	e *echo.Echo,
	linkHandler handler.ILinkHandler,
	postHandler handler.IPostHandler,
	authHandler handler.IAuthHandler,
	rssHandler handler.IRssHandler,
	modelHandler handler.IModelHandler,
	am *middleware.AuthMiddleware,
) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// 路由注册
	router.RegisterLinkRoutes(v1, linkHandler)
	router.RegisterPostRoutes(v1, am, postHandler)
	router.RegisterAuthRoutes(v1, authHandler)
	router.RegisterRssRoutes(v1, rssHandler)
	router.RegisterModelRoutes(v1, modelHandler)
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

func runServerLifecycle(lc fx.Lifecycle, app *echo.Echo, cfg *config.Config, log logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				log.Info("Server is starting")
				if err := app.Start(cfg.Server.GetAddr()); err != nil {
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

			if err := app.Shutdown(shutdownCtx); err != nil {
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
		database.NewDB,
		cache.NewCacheClient,
		validatorx.NewValidator,
		providerEchoApp,

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
		service.NewModelService,

		// handlers
		handler.NewAuthHandler,
		handler.NewPostHandler,
		handler.NewLinkHandler,
		handler.NewRssHandler,
		handler.NewModelHandler,
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
