// Package repository
package repository

import (
	"context"
	"errors"
	"strings"

	"blog-server/internal/entity"
	"blog-server/pkg/errs"

	"gorm.io/gorm"
)

type ILinkRepo interface {
	GetAllLinks(ctx context.Context, db *gorm.DB) ([]*entity.Link, error)
	GetAllEnabledLinks(ctx context.Context, db *gorm.DB) ([]*entity.Link, error)
	CreateLink(ctx context.Context, db *gorm.DB, link *entity.Link) error
	UpdateLinkStatusBatch(ctx context.Context, db *gorm.DB, update map[uint]entity.LinkStatus) error
}

type linkRepo struct{}

// UpdateLinkStatusBatch implements [ILinkRepo].
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

func NewLinkRepo() ILinkRepo {
	return &linkRepo{}
}

func (r *linkRepo) GetAllLinks(ctx context.Context, db *gorm.DB) ([]*entity.Link, error) {
	links, err := gorm.G[*entity.Link](db).Find(ctx)
	if err != nil {
		return nil, errs.New(errs.CodeDatabaseError, "database error", err)
	}
	return links, nil
}

func (r *linkRepo) GetAllEnabledLinks(ctx context.Context, db *gorm.DB) ([]*entity.Link, error) {
	links, err := gorm.G[*entity.Link](db).Where("enabled = ?", true).Find(ctx)
	if err != nil {
		return nil, errs.New(errs.CodeDatabaseError, "database error", err)
	}
	return links, nil
}

func (r *linkRepo) CreateLink(ctx context.Context, db *gorm.DB, link *entity.Link) error {
	if err := gorm.G[entity.Link](db).Create(ctx, link); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.New(errs.CodeConflict, "link already exists", err)
		}
		return errs.New(errs.CodeDatabaseError, "database error", err)
	}
	return nil
}
