package jobs

import (
	"context"
	"log"
	"time"

	"blog-server/service"
)

func StartCheckLinkStatusJob(ctx context.Context, svc service.ILinkService) {
	for {
		now := time.Now()
		next := now.Truncate(time.Hour).Add(time.Hour)
		wait := time.Until(next)

		select {
		case <-time.After(wait):
			if err := svc.CheckLinkStatus(ctx); err != nil {
				log.Printf("[check_links_status] check error: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}
