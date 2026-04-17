package main

import (
	"context"
	"time"

	"blog-server/cache"
	"blog-server/config"
	"blog-server/datastore"
	"blog-server/handler"
	"blog-server/logger"
	"blog-server/middleware"
	"blog-server/pkg/validatorx"
	"blog-server/repository"
	"blog-server/scheduler"
	"blog-server/service"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
)

// main initializes the application using Uber Fx and starts the dependency
// injection container.
//
// It registers configuration, logging, datastore, HTTP server, and lifecycle
// hooks required to run the service.
func main() {
	app := fx.New(
		fx.Options(
			config.Module(),
			logger.Module(),
			cache.Module(),
			datastore.Module(),
			repository.Module(),
			service.Module(),
			handler.Module(),
			scheduler.Module(),
		),
		fx.Provide(
			validatorx.NewValidator,
			middleware.NewAuthMiddleware,
			middleware.NewRoleMiddleware,
			providerFiberApp,
		),
		fx.Invoke(
			runServerLifecycle,
		),
	)
	app.Run()
}

// providerFiberApp constructs and configures the Fiber HTTP server instance.
//
// It applies base server configuration such as proxy handling and request limits.
func providerFiberApp(cfg *config.Config, log logger.Logger) *fiber.App {
	fiberCfg := fiber.Config{
		ErrorHandler: handler.ErrorHandler(log, cfg),
		TrustProxy:   true,
		ProxyHeader:  fiber.HeaderXForwardedFor,
		BodyLimit:    10 * 1024 * 1024,
	}

	app := fiber.New(fiberCfg)
	app.Use(middleware.RequestLogger(log, cfg))
	return app
}

// runServerLifecycle registers Fiber server startup and graceful shutdown
// hooks into the Fx application lifecycle.
//
// OnStart:
//   - Starts the HTTP server in a separate goroutine
//
// OnStop:
//   - Performs graceful shutdown using configured timeout
//   - Ensures in-flight requests are completed before exit
func runServerLifecycle(lc fx.Lifecycle, app *fiber.App, cfg *config.Config, log logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Info("Server is starting")
				if err := app.Listen(cfg.Server.Addr()); err != nil {
					log.Error("Server startup failed", logger.Err(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			timeout := cfg.Server.GracefulShutdown
			if timeout <= 0 {
				timeout = 5 * time.Second
			}
			shutdownCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			if err := app.ShutdownWithContext(shutdownCtx); err != nil {
				log.Error("Server shutdown failed", logger.Err(err))
			} else {
				log.Info("Server has been shut down successfully.")
			}
			return nil
		},
	})
}
