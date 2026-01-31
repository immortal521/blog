// Package scheduler
package scheduler

import (
	"context"

	"blog-server/scheduler/jobs"
	"blog-server/service"
)

type Scheduler struct {
	postService service.IPostService
	linkService service.ILinkService
}

func NewScheduler(postService service.IPostService, linkService service.ILinkService) *Scheduler {
	return &Scheduler{postService, linkService}
}

func (s *Scheduler) Start(ctx context.Context) {
	go jobs.StartViewFlushJob(ctx, s.postService)
	go jobs.StartCheckLinkStatusJob(ctx, s.linkService)
}
