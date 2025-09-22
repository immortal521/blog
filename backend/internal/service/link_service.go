package service

import (
	"blog-server/internal/database"
	"blog-server/internal/dto/request"
	"blog-server/internal/entity"
	"blog-server/internal/repo"
	"context"
)

type ILinkService interface {
	GetLinks(ctx context.Context) ([]*entity.Link, error)
	CreateLink(ctx context.Context, dto *request.CreateLinkReq) error
}

type linkService struct {
	db       database.DB
	linkRepo repo.LinkRepo
}

func NewLinkService(db database.DB, linkRepo repo.LinkRepo) ILinkService {
	return &linkService{db: db, linkRepo: linkRepo}
}

func (s *linkService) GetLinks(ctx context.Context) ([]*entity.Link, error) {
	links, err := s.linkRepo.GetAllLinks(ctx, s.db.Conn())
	if err != nil {
		return nil, err
	}
	return links, nil
}

func (s *linkService) CreateLink(ctx context.Context, dto *request.CreateLinkReq) error {
	link := entity.Link{
		Name:        dto.Name,
		Description: dto.Description,
		Avatar:      dto.Avatar,
		URL:         dto.URL,
	}
	err := s.linkRepo.CreateLink(ctx, s.db.Conn(), &link)
	if err != nil {
		return err
	}
	return nil
}
