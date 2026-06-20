package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"blog-server/authz"
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

	"github.com/labstack/echo/v5"
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
			authz.Module(),
			service.Module(),
			handler.Module(),
			scheduler.Module(),
		),
		fx.Provide(
			validatorx.NewValidator,
			middleware.NewAuthMiddleware,
			providerEchoApp,
		),
		fx.Invoke(
			runServerLifecycle,
		),
	)
	app.Run()
}

func providerEchoApp(cfg *config.Config, log logger.Logger) *echo.Echo {
	echoCfg := echo.Config{
		HTTPErrorHandler: handler.ErrorHandler(cfg, log),
		IPExtractor:      echo.ExtractIPFromXFFHeader(),
	}
	app := echo.NewWithConfig(echoCfg)
	app.Use(middleware.RequestLogger(cfg, log))
	app.Use(middleware.BodyLimit(10 * 10 * 1024))
	return app
}

// runServerLifecycle registers Echo server startup and graceful shutdown
// hooks into the Fx application lifecycle.
//
// OnStart:
//   - Starts the HTTP server in a separate goroutine
//
// OnStop:
//   - Performs graceful shutdown using configured timeout
//   - Ensures in-flight requests are completed before exit
func runServerLifecycle(lc fx.Lifecycle, app *echo.Echo, cfg *config.Config, log logger.Logger) {
	srv := &http.Server{
		Addr:    cfg.Server.Addr(),
		Handler: app,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Info("Server is starting")
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
			shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			log.Info("Server is shutting down")
			if err := srv.Shutdown(shutdownCtx); err != nil {
				log.Error("Server shutdown failed", logger.Err(err))
				return err
			}
			log.Info("Server has been shut down successfully")
			return nil
		},
	})
}
