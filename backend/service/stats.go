package service

import (
	"context"

	"blog-server/database"
	"blog-server/repository"
	"blog-server/response"
)

type IStatsService interface {
	GetDashboardStats(ctx context.Context) (*response.DashboardStatsRes, error)
}

type statsService struct {
	db       database.Database
	postRepo repository.IPostRepo
}

func (s *statsService) GetDashboardStats(ctx context.Context) (*response.DashboardStatsRes, error) {
	postCount, err := s.postRepo.GetPostCount(ctx, s.db)
	if err != nil {
		return nil, err
	}
	return &response.DashboardStatsRes{
		TotalPosts: postCount,
	}, nil
}

func NewStatsService(db database.Database, postRepo repository.IPostRepo) IStatsService {
	return &statsService{
		db,
		postRepo,
	}
}
