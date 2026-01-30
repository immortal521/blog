// Package repository provides data access layer for the blog system
package repository

import (
	"context"
	"errors"
	"strings"

	"blog-server/internal/entity"
	"blog-server/internal/response"
	"blog-server/pkg/errs"

	"gorm.io/gorm"
)

// ILinkRepo defines the interface for link data access operations
type ILinkRepo interface {
	GetAllLinks(ctx context.Context, db *gorm.DB) ([]*entity.Link, error)
	GetAllEnabledLinks(ctx context.Context, db *gorm.DB) ([]*entity.Link, error)
	CreateLink(ctx context.Context, db *gorm.DB, link *entity.Link) error
	UpdateLinkStatusBatch(ctx context.Context, db *gorm.DB, updates map[uint]entity.LinkStatus) error
	GetOverview(ctx context.Context, db *gorm.DB) (*response.LinkOverview, error)
}

type linkRepo struct{}

// NewLinkRepo creates a new link repository instance
func NewLinkRepo() ILinkRepo {
	return &linkRepo{}
}

// GetAllLinks retrieves all links from the database
func (r *linkRepo) GetAllLinks(ctx context.Context, db *gorm.DB) ([]*entity.Link, error) {
	links, err := gorm.G[*entity.Link](db).Find(ctx)
	if err != nil {
		return nil, errs.New(errs.CodeDatabaseError, "database error", err)
	}
	return links, nil
}

// GetAllEnabledLinks retrieves all enabled links from the database
func (r *linkRepo) GetAllEnabledLinks(ctx context.Context, db *gorm.DB) ([]*entity.Link, error) {
	links, err := gorm.G[*entity.Link](db).Where("enabled = ?", true).Find(ctx)
	if err != nil {
		return nil, errs.New(errs.CodeDatabaseError, "database error", err)
	}
	return links, nil
}

// CreateLink creates a new link in the database
func (r *linkRepo) CreateLink(ctx context.Context, db *gorm.DB, link *entity.Link) error {
	if err := gorm.G[entity.Link](db).Create(ctx, link); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.New(errs.CodeConflict, "link already exists", err)
		}
		return errs.New(errs.CodeDatabaseError, "database error", err)
	}
	return nil
}

// UpdateLinkStatusBatch batch updates the status of multiple links
func (r *linkRepo) UpdateLinkStatusBatch(ctx context.Context, db *gorm.DB, updates map[uint]entity.LinkStatus) error {
	if len(updates) == 0 {
		return nil
	}

	var caseBuilder strings.Builder
	var idArgs []any
	var ids []any

	caseBuilder.WriteString("CASE id ")
	for id, status := range updates {
		caseBuilder.WriteString("WHEN ? THEN ? ")
		idArgs = append(idArgs, id, status)
		ids = append(ids, id)
	}
	caseBuilder.WriteString("ELSE status END")

	err := db.Model(&entity.Link{}).
		Where("id IN (?)", ids).
		UpdateColumn("status", gorm.Expr(caseBuilder.String(), idArgs...)).Error
	if err != nil {
		return errs.New(errs.CodeDatabaseError, "database error", err)
	}
	return nil
}

// GetOverview retrieves link statistics including total, normal, abnormal, and pending counts
func (r *linkRepo) GetOverview(ctx context.Context, db *gorm.DB) (*response.LinkOverview, error) {
	type result struct {
		Total   int
		Normal  int
		Enabled int
	}

	var res result
	err := db.Table("links").
		Select(
			"count(id) as total, "+
				"sum(case when status = ? then 1 else 0 end) as normal, "+
				"sum(case when enabled = ? then 1 else 0 end) as enabled",
			entity.LinkNormal, true,
		).Scan(&res).Error
	if err != nil {
		return nil, errs.New(errs.CodeDatabaseError, "database error", err)
	}

	return &response.LinkOverview{
		Total:    res.Total,
		Normal:   res.Normal,
		Abnormal: res.Total - res.Normal,
		Pending:  res.Total - res.Enabled,
	}, nil
}
