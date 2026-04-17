// Package scheduler
package scheduler

import (
	"context"

	"blog-server/scheduler/jobs"
	"blog-server/service"
)

type Scheduler struct {
	postService service.PostService
	linkService service.LinkService
}

func NewScheduler(postService service.PostService, linkService service.LinkService) *Scheduler {
	return &Scheduler{postService, linkService}
}

func (s *Scheduler) Start(ctx context.Context) {
	go jobs.StartViewFlushJob(ctx, s.postService)
	go jobs.StartCheckLinkStatusJob(ctx, s.linkService)
}
