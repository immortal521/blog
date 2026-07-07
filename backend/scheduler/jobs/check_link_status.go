package jobs

import (
	"context"
	"time"

	"blog-server/logger"
	"blog-server/service"
)

func StartCheckLinkStatusJob(ctx context.Context, svc service.LinkService, log logger.Logger) {
	for {
		now := time.Now()
		next := now.Truncate(time.Hour).Add(time.Hour)
		wait := time.Until(next)

		select {
		case <-time.After(wait):
			if err := svc.CheckLinkStatus(ctx); err != nil {
				log.Error("check link status failed",
					logger.String("module", "scheduler"),
					logger.String("job", "check_link_status"),
					logger.Err(err),
				)
			}
		case <-ctx.Done():
			return
		}
	}
}
