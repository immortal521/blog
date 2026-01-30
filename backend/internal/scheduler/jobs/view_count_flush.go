// Package jobs provider scheduler jobs
package jobs

import (
	"context"
	"log"
	"time"

	"blog-server/internal/service"
)

func StartViewFlushJob(ctx context.Context, svc service.IPostService) {
	for {
		now := time.Now()
		next := now.Truncate(time.Hour).Add(time.Hour)
		wait := time.Until(next)

		select {
		case <-time.After(wait):
			if err := svc.FlushViewCountToDB(ctx); err != nil {
				log.Printf("[view_flush] flush error: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}
