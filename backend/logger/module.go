package logger

import "go.uber.org/fx"

// Module registers the logger implementation into the Fx dependency injection container.
//
// It provides a single Logger instance via NewLogger for application-wide use.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			NewLogger,
		),
	)
}
