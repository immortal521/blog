// Package jobs provider scheduler jobs
package jobs

import (
	"context"
	"time"

	"blog-server/logger"
	"blog-server/service"
)

func StartViewFlushJob(ctx context.Context, svc service.PostService, log logger.Logger) {
	for {
		now := time.Now()
		next := now.Truncate(time.Hour).Add(time.Hour)
		wait := time.Until(next)

		select {
		case <-time.After(wait):
			if err := svc.FlushViewCountToDB(ctx); err != nil {
				log.Error("check link status failed",
					logger.String("module", "scheduler"),
					logger.String("job", "view_count_flush"),
					logger.Err(err),
				)
			}
		case <-ctx.Done():
			return
		}
	}
}
