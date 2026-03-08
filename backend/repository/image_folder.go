package repository

import (
	"context"
	"errors"
	"time"

	"blog-server/database"
	"blog-server/entity"
	"blog-server/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IImageFolderRepo interface {
	Create(ctx context.Context, db database.DB, folder *entity.ImageFolder) error
	GetByID(ctx context.Context, db database.DB, id uuid.UUID) (*entity.ImageFolder, error)
	Exists(ctx context.Context, db database.DB, id uuid.UUID) (bool, error)
	ListByParent(ctx context.Context, db database.DB, parentID *uuid.UUID, limit, offset int) ([]entity.ImageFolder, error)
	ExistsBySameNameInParent(ctx context.Context, db database.DB, parentID *uuid.UUID, name string, excludeID *uuid.UUID) (bool, error)
	CountChildren(ctx context.Context, db database.DB, parentID *uuid.UUID) (int64, error)
	Rename(ctx context.Context, db database.DB, id uuid.UUID, newName string) error
	Move(ctx context.Context, db database.DB, id uuid.UUID, newParentID *uuid.UUID) error
	SoftDelete(ctx context.Context, db database.DB, id uuid.UUID) error
}

type imageFolderRepo struct{}

func (i *imageFolderRepo) CountChildren(ctx context.Context, db database.DB, parentID *uuid.UUID) (int64, error) {
	gdb := unwrapDB(db)

	query := gorm.G[entity.ImageFolder](gdb).
		Where("deleted_at IS NULL")

	if parentID != nil {
		query = query.Where("parent_id = ?", parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	count, err := query.Count(ctx, "id")
	if err != nil {
		return 0, errs.New(errs.CodeDatabaseError, "database error", err)
	}

	return count, nil
}

func (i *imageFolderRepo) Create(ctx context.Context, db database.DB, folder *entity.ImageFolder) error {
	gdb := unwrapDB(db)

	if folder.ID == uuid.Nil {
		folder.ID = uuid.New()
	}

	now := time.Now()

	if folder.CreatedAt.IsZero() {
		folder.CreatedAt = now
	}

	folder.UpdatedAt = now

	err := gorm.G[entity.ImageFolder](gdb).Create(ctx, folder)
	if err != nil {
		return errs.New(errs.CodeDatabaseError, "database error", err)
	}

	return nil
}

func (i *imageFolderRepo) Exists(ctx context.Context, db database.DB, id uuid.UUID) (bool, error) {
	gdb := unwrapDB(db)

	count, err := gorm.G[entity.ImageFolder](gdb).
		Where("id = ? AND deleted_at IS NULL", id).
		Count(ctx, "id")
	if err != nil {
		return false, errs.New(errs.CodeDatabaseError, "database error", err)
	}

	return count > 0, nil
}

func (i *imageFolderRepo) ExistsBySameNameInParent(
	ctx context.Context,
	db database.DB,
	parentID *uuid.UUID,
	name string,
	excludeID *uuid.UUID,
) (bool, error) {
	gdb := unwrapDB(db)

	query := gorm.G[entity.ImageFolder](gdb).
		Where("deleted_at IS NULL").
		Where("name = ?", name)

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}

	count, err := query.Count(ctx, "id")
	if err != nil {
		return false, errs.New(errs.CodeDatabaseError, "database error", err)
	}

	return count > 0, nil
}

func (i *imageFolderRepo) GetByID(ctx context.Context, db database.DB, id uuid.UUID) (*entity.ImageFolder, error) {
	gdb := unwrapDB(db)

	folder, err := gorm.G[*entity.ImageFolder](gdb).
		Where("id = ? AND deleted_at IS NULL", id).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, errs.New(errs.CodeDatabaseError, "database error", err)
	}

	return folder, nil
}

func (i *imageFolderRepo) ListByParent(
	ctx context.Context,
	db database.DB,
	parentID *uuid.UUID,
	limit int,
	offset int,
) ([]entity.ImageFolder, error) {
	gdb := unwrapDB(db)

	query := gorm.G[entity.ImageFolder](gdb).Where("deleted_at IS NULL")

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", parentID)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	result, err := query.Order("name ASC").Find(ctx)
	if err != nil {
		return nil, errs.New(errs.CodeDatabaseError, "database error", err)
	}

	return result, nil
}

func (i *imageFolderRepo) Move(
	ctx context.Context,
	db database.DB,
	id uuid.UUID,
	newParentID *uuid.UUID,
) error {
	gdb := unwrapDB(db)

	_, err := gorm.G[entity.ImageFolder](gdb).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(ctx, entity.ImageFolder{
			ParentID:  newParentID,
			UpdatedAt: time.Now(),
		})
	if err != nil {
		return errs.New(errs.CodeDatabaseError, "database error", err)
	}

	return nil
}

func (i *imageFolderRepo) Rename(
	ctx context.Context,
	db database.DB,
	id uuid.UUID,
	newName string,
) error {
	gdb := unwrapDB(db)

	_, err := gorm.G[entity.ImageFolder](gdb).
		Where("id = ?", id).
		Updates(ctx, entity.ImageFolder{
			Name:      newName,
			UpdatedAt: time.Now(),
		})
	if err != nil {
		return errs.New(errs.CodeDatabaseError, "database error", err)
	}

	return nil
}

func (i *imageFolderRepo) SoftDelete(
	ctx context.Context,
	db database.DB,
	id uuid.UUID,
) error {
	gdb := unwrapDB(db)

	now := time.Now()

	_, err := gorm.G[entity.ImageFolder](gdb).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(ctx, entity.ImageFolder{
			DeletedAt: &now,
			UpdatedAt: now,
		})
	if err != nil {
		return errs.New(errs.CodeDatabaseError, "database error", err)
	}

	return nil
}

func NewImageFolderRepo() IImageFolderRepo {
	return &imageFolderRepo{}
}
