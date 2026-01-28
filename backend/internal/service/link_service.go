package service

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"blog-server/internal/database"
	"blog-server/internal/entity"
	"blog-server/internal/repository"
	"blog-server/internal/request"
)

type ILinkService interface {
	GetLinks(ctx context.Context) ([]*entity.Link, error)
	CreateLink(ctx context.Context, dto *request.CreateLinkReq) error
	CheckLinkStatus(ctx context.Context) error
}

type linkService struct {
	db       database.DB
	linkRepo repository.ILinkRepo
}

func NewLinkService(db database.DB, linkRepo repository.ILinkRepo) ILinkService {
	return &linkService{db: db, linkRepo: linkRepo}
}

func (s *linkService) GetLinks(ctx context.Context) ([]*entity.Link, error) {
	return s.linkRepo.GetAllEnabledLinks(ctx, s.db.Conn())
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

func (s linkService) CheckLinkStatus(ctx context.Context) error {
	links, err := s.linkRepo.GetAllLinks(ctx, s.db.Conn())
	if err != nil {
		return err
	}

	updates := make(map[uint]entity.LinkStatus)
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	sem := make(chan struct{}, 10)

	for _, link := range links {
		wg.Add(1)
		sem <- struct{}{}

		go func(l *entity.Link) {
			defer wg.Done()
			defer func() { <-sem }()

			currentStatus := l.Status
			status := entity.LinkAbnormal

			if strings.HasPrefix(l.URL, "https://") {
				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Get(l.URL)
				if err == nil && resp.StatusCode == http.StatusOK {
					status = entity.LinkNormal
				}
				if resp != nil {
					resp.Body.Close()
				}
			}

			if status != currentStatus {
				mu.Lock()
				link.Status = status
				updates[l.ID] = status
				mu.Unlock()
			}
		}(link)
	}

	wg.Wait()
	err = s.linkRepo.UpdateLinkStatusBatch(ctx, s.db.Conn(), updates)
	if err != nil {
		return err
	}

	return nil
}
