package repo

import (
	"blog-server/internal/entity"
	"blog-server/pkg/errs"
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type LinkRepo interface {
	GetAllLinks(db *gorm.DB) ([]entity.Link, error)
	CreateLink(db *gorm.DB, link *entity.Link) error
}

type linkRepo struct{}

func NewLinkRepo() LinkRepo {
	return &linkRepo{}
}

func (r *linkRepo) GetAllLinks(db *gorm.DB) ([]entity.Link, error) {
	var links []entity.Link
	err := db.Find(&links).Where("enabled = ?", true).Error
	if err != nil {
		return nil, err
	}
	return links, nil
}

func (r *linkRepo) CreateLink(db *gorm.DB, link *entity.Link) error {
	err := db.Create(link).Error
	log.Error(errors.Is(err, gorm.ErrRecordNotFound))
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.BadRequest("url is duplicated")
		}
		return err
	}
	return nil
}
