// Package repo
package repo

import (
	"blog-server/internal/entity"
	"blog-server/pkg/errs"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type LinkRepo interface {
	GetAllLinks(ctx context.Context, db *gorm.DB) ([]entity.Link, error)
	CreateLink(ctx context.Context, db *gorm.DB, link *entity.Link) error
}

type linkRepo struct{}

func NewLinkRepo() LinkRepo {
	return &linkRepo{}
}

func (r *linkRepo) GetAllLinks(ctx context.Context, db *gorm.DB) ([]entity.Link, error) {
	links, err := gorm.G[entity.Link](db).Where("enabled = ?", true).Find(ctx)
	if err != nil {
		return nil, err
	}
	return links, nil
}

func (r *linkRepo) CreateLink(ctx context.Context, db *gorm.DB, link *entity.Link) error {
	err := gorm.G[entity.Link](db).Create(ctx, link)
	log.Error(errors.Is(err, gorm.ErrRecordNotFound))
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.BadRequest("url is duplicated")
		}
		return err
	}
	return nil
}
