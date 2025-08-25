package service

import (
	"blog-server/internal/database"
	"blog-server/internal/dto"
	"blog-server/internal/entity"
	"blog-server/internal/repo"
	"context"
)

type ILinkService interface {
	GetLinks(ctx context.Context) ([]entity.Link, error)
	CreateLink(ctx context.Context, dto *dto.LinkCreateReq) error
}

type linkService struct {
	db       database.DB
	linkRepo repo.LinkRepo
}

func NewLinkService(db database.DB, linkRepo repo.LinkRepo) ILinkService {
	return &linkService{db: db, linkRepo: linkRepo}
}

func (s *linkService) GetLinks(ctx context.Context) ([]entity.Link, error) {
	links, err := s.linkRepo.GetAllLinks(s.db.Conn(ctx))
	if err != nil {
		return nil, err
	}
	return links, nil
}

func (s *linkService) CreateLink(ctx context.Context, dto *dto.LinkCreateReq) error {
	link := entity.Link{
		Name:        dto.Name,
		Description: dto.Description,
		Avatar:      dto.Avatar,
		URL:         dto.Url,
	}
	err := s.linkRepo.CreateLink(s.db.Conn(ctx), &link)
	if err != nil {
		return err
	}
	return nil
}
