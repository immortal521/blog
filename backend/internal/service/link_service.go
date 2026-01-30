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
	"blog-server/internal/response"
)

// ILinkService defines the interface for link business logic operations
type ILinkService interface {
	GetLinks(ctx context.Context) ([]*entity.Link, error)
	CreateLink(ctx context.Context, dto *request.CreateLinkReq) error
	CheckLinkStatus(ctx context.Context) error
	GetOverview(ctx context.Context) (*response.LinkOverview, error)
}

type linkService struct {
	db       database.DB
	linkRepo repository.ILinkRepo
}

// NewLinkService creates a new link service instance
func NewLinkService(db database.DB, linkRepo repository.ILinkRepo) ILinkService {
	return &linkService{db: db, linkRepo: linkRepo}
}

// GetLinks retrieves all enabled links
func (s *linkService) GetLinks(ctx context.Context) ([]*entity.Link, error) {
	return s.linkRepo.GetAllEnabledLinks(ctx, s.db.Conn())
}

// GetOverview retrieves link statistics
func (s *linkService) GetOverview(ctx context.Context) (*response.LinkOverview, error) {
	return s.linkRepo.GetOverview(ctx, s.db.Conn())
}

// CreateLink creates a new link
func (s *linkService) CreateLink(ctx context.Context, dto *request.CreateLinkReq) error {
	link := entity.Link{
		Name:        dto.Name,
		Description: dto.Description,
		Avatar:      dto.Avatar,
		URL:         dto.URL,
	}
	return s.linkRepo.CreateLink(ctx, s.db.Conn(), &link)
}

// CheckLinkStatus checks the status of all links by making HTTP requests
// Updates the link status in the database if it has changed
func (s linkService) CheckLinkStatus(ctx context.Context) error {
	links, err := s.linkRepo.GetAllLinks(ctx, s.db.Conn())
	if err != nil {
		return err
	}

	updates := make(map[uint]entity.LinkStatus)
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	sem := make(chan struct{}, 10) // Limit to 10 concurrent requests

	for _, link := range links {
		wg.Add(1)
		sem <- struct{}{}

		go func(l *entity.Link) {
			defer wg.Done()
			defer func() { <-sem }()

			currentStatus := l.Status
			status := entity.LinkAbnormal

			// Check HTTPS links only
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

			// Update status if changed
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
