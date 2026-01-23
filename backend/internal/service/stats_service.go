package service

import (
	"context"

	"blog-server/internal/database"
	"blog-server/internal/repository"
	"blog-server/internal/response"
)

type IStatsService interface {
	GetDashboardStats(ctx context.Context) (*response.DashboardStatsRes, error)
}

type statsService struct {
	db       database.DB
	postRepo repository.IPostRepo
}

func (s *statsService) GetDashboardStats(ctx context.Context) (*response.DashboardStatsRes, error) {
	postCount, err := s.postRepo.GetPostCount(ctx, s.db.Conn())
	if err != nil {
		return nil, err
	}
	return &response.DashboardStatsRes{
		TotalPosts: postCount,
	}, nil
}

func NewStatsService(db database.DB, postRepo repository.IPostRepo) IStatsService {
	return &statsService{
		db,
		postRepo,
	}
}
