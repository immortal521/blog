package repository

import (
	"context"
	"errors"
	"time"

	"blog-server/entity"
	"blog-server/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IImageRepo interface {
	Create(ctx context.Context, db *gorm.DB, image *entity.Image) error
	GetByID(ctx context.Context, db *gorm.DB, id uuid.UUID) (*entity.Image, error)
	ListByFolder(ctx context.Context, db *gorm.DB, folderID *uuid.UUID, limit, offset int) ([]entity.Image, error)
	Move(ctx context.Context, db *gorm.DB, id uuid.UUID, newFolderID *uuid.UUID) error
	SoftDelete(ctx context.Context, db *gorm.DB, id uuid.UUID) error
	CountByFolder(ctx context.Context, db *gorm.DB, folderID *uuid.UUID) (int64, error)
	ExistsBySameNameInFolder(ctx context.Context, db *gorm.DB, folderID *uuid.UUID, name string, excludeID *uuid.UUID) (bool, error)
	GetBySha256(ctx context.Context, db *gorm.DB, sha256 string) (*entity.Image, error)
}

type imageRepo struct{}

func (i *imageRepo) GetBySha256(ctx context.Context, db *gorm.DB, sha256 string) (*entity.Image, error) {
	image, err := gorm.G[*entity.Image](db).
		Where("sha256 = ? AND deleted_at IS NULL", sha256).
		First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errs.New(errs.CodeDatabaseError, "database error", err)
	}
	return image, err
}

func (i *imageRepo) CountByFolder(ctx context.Context, db *gorm.DB, folderID *uuid.UUID) (int64, error) {
	query := gorm.G[entity.Image](db).Where("deleted_at IS NULL")
	if folderID == nil {
		query = query.Where("folder_id IS NULL")
	} else {
		query = query.Where("folder_id = ?", *folderID)
	}

	count, err := query.Count(ctx, "id")
	if err != nil {
		return 0, errs.New(errs.CodeDatabaseError, "database error", err)
	}
	return count, nil
}

func (i *imageRepo) Create(ctx context.Context, db *gorm.DB, image *entity.Image) error {
	if image.ID == uuid.Nil {
		image.ID = uuid.New()
	}
	now := time.Now()
	if image.CreatedAt.IsZero() {
		image.CreatedAt = time.Now()
	}
	image.UpdatedAt = now
	err := gorm.G[entity.Image](db).Create(ctx, image)
	if err != nil {
		return errs.New(errs.CodeDatabaseError, "database error", err)
	}
	return nil
}

func (i *imageRepo) GetByID(ctx context.Context, db *gorm.DB, id uuid.UUID) (*entity.Image, error) {
	image, err := gorm.G[*entity.Image](db).
		Where("id = ? AND deleted_at IS NULL", id).
		First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errs.New(errs.CodeDatabaseError, "database error", err)
	}

	return image, nil
}

func (i *imageRepo) ListByFolder(ctx context.Context, db *gorm.DB, folderID *uuid.UUID, limit int, offset int) ([]entity.Image, error) {
	query := gorm.G[entity.Image](db).Where("deleted_at IS NULL")

	if folderID == nil {
		query = query.Where("folder_id IS NULL")
	} else {
		query = query.Where("folder_id = ?", folderID)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	images, err := query.Order("origin_name DESC").Find(ctx)
	if err != nil {
		return nil, errs.New(errs.CodeDatabaseError, "database error", err)
	}

	return images, nil
}

func (i *imageRepo) Move(ctx context.Context, db *gorm.DB, id uuid.UUID, newFolderID *uuid.UUID) error {
	_, err := gorm.G[entity.Image](db).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(ctx, entity.Image{
			FolderID:  newFolderID,
			UpdatedAt: time.Now(),
		})
	if err != nil {
		return errs.New(errs.CodeDatabaseError, "database error", err)
	}
	return nil
}

func (i *imageRepo) ExistsBySameNameInFolder(
	ctx context.Context,
	db *gorm.DB,
	folderID *uuid.UUID,
	name string,
	excludeID *uuid.UUID,
) (bool, error) {
	q := gorm.G[entity.Image](db).
		Where("deleted_at IS NULL").
		Where("origin_name = ?", name)

	if folderID == nil {
		q = q.Where("folder_id IS NULL")
	} else {
		q = q.Where("folder_id = ?", *folderID)
	}

	if excludeID != nil {
		q = q.Where("id <> ?", *excludeID)
	}

	count, err := q.Count(ctx, "id")
	if err != nil {
		return false, errs.New(errs.CodeDatabaseError, "database error", err)
	}
	return count > 0, nil
}

func (i *imageRepo) SoftDelete(ctx context.Context, db *gorm.DB, id uuid.UUID) error {
	now := time.Now()
	_, err := gorm.G[entity.Image](db).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(ctx, entity.Image{
			DeletedAt: &now,
			UpdatedAt: now,
		})
	if err != nil {
		return errs.New(errs.CodeDatabaseError, "database error", err)
	}
	return nil
}

func NewImageRepo() IImageRepo {
	return &imageRepo{}
}
