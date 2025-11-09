package service

import (
	"blog-server/internal/database"
	"blog-server/internal/entity"
	"blog-server/internal/repo"
	"blog-server/internal/request"
	"context"
)

type ILinkService interface {
	GetLinks(ctx context.Context) ([]*entity.Link, error)
	CreateLink(ctx context.Context, dto *request.CreateLinkReq) error
}

type linkService struct {
	db       database.DB
	linkRepo repo.ILinkRepo
}

func NewLinkService(db database.DB, linkRepo repo.ILinkRepo) ILinkService {
	return &linkService{db: db, linkRepo: linkRepo}
}

func (s *linkService) GetLinks(ctx context.Context) ([]*entity.Link, error) {
	return s.linkRepo.GetAllLinks(ctx, s.db.Conn())
}

func (s *linkService) CreateLink(ctx context.Context, dto *request.CreateLinkReq) error {
	link := entity.Link{
		Name:        dto.Name,
		Description: dto.Description,
		Avatar:      dto.Avatar,
		URL:         dto.URL,
	}
	return s.linkRepo.CreateLink(ctx, s.db.Conn(), &link)
}
