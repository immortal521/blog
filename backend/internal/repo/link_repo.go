// Package repo
package repo

import (
	"blog-server/internal/entity"
	"blog-server/pkg/errs"
	"context"
	"errors"

	"gorm.io/gorm"
)

type LinkRepo interface {
	GetAllLinks(ctx context.Context, db *gorm.DB) ([]*entity.Link, error)
	CreateLink(ctx context.Context, db *gorm.DB, link *entity.Link) error
}

type linkRepo struct{}

func NewLinkRepo() LinkRepo {
	return &linkRepo{}
}

func (r *linkRepo) GetAllLinks(ctx context.Context, db *gorm.DB) ([]*entity.Link, error) {
	links, err := gorm.G[*entity.Link](db).Where("enabled = ?", true).Find(ctx)
	if err != nil {
		return nil, err
	}
	return links, nil
}

func (r *linkRepo) CreateLink(ctx context.Context, db *gorm.DB, link *entity.Link) error {
	err := gorm.G[entity.Link](db).Create(ctx, link)
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return errs.ErrDuplicateURL
	}
	return err
}
