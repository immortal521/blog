// Package scheduler
package scheduler

import (
	"context"

	"blog-server/internal/scheduler/jobs"
	"blog-server/internal/service"
)

type Scheduler struct {
	postService service.IPostService
}

func NewScheduler(postService service.IPostService) *Scheduler {
	return &Scheduler{postService: postService}
}

func (s *Scheduler) Start(ctx context.Context) {
	go jobs.StartViewFlushJob(ctx, s.postService)
}
