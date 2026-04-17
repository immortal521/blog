package scheduler

import (
	"context"

	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"scheduler",
		fx.Provide(
			NewScheduler,
		),
		fx.Invoke(
			runJobsLifecycle,
		),
	)
}

func runJobsLifecycle(lc fx.Lifecycle, scheduler *Scheduler) {
	var cancel context.CancelFunc

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			jobCtx, c := context.WithCancel(context.Background())
			cancel = c

			go scheduler.Start(jobCtx)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if cancel != nil {
				cancel()
			}
			return nil
		},
	})
}
