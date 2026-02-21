package service

import (
	"context"
	"strings"

	"blog-server/database"
	"blog-server/entity"
	"blog-server/errs"
	"blog-server/repository"

	"github.com/google/uuid"
)

type IImageFolderService interface {
	Create(ctx context.Context, name string, parentID *uuid.UUID) (*entity.ImageFolder, error)
	Get(ctx context.Context, id uuid.UUID) (*entity.ImageFolder, error)
	List(ctx context.Context, parentID *uuid.UUID, limit, offset int) ([]entity.ImageFolder, error)
	Rename(ctx context.Context, id uuid.UUID, newName string) error
	Move(ctx context.Context, id uuid.UUID, targetParentID *uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type imageFolderService struct {
	db              database.DB
	imageFolderRepo repository.IImageFolderRepo
}

func normalizeName(name string) string {
	return strings.TrimSpace(name)
}

func (s *imageFolderService) ensureFolderExists(ctx context.Context, id uuid.UUID) (*entity.ImageFolder, error) {
	folder, err := s.imageFolderRepo.GetByID(ctx, s.db.Conn(), id)
	if err != nil {
		return nil, err
	}
	if folder == nil {
		return nil, errs.New(errs.CodeNoContent, "no such folder", nil)
	}
	return folder, nil
}

func (s *imageFolderService) ensureParentExists(ctx context.Context, parentID *uuid.UUID) error {
	if parentID == nil {
		return nil
	}
	ok, err := s.imageFolderRepo.Exists(ctx, s.db.Conn(), *parentID)
	if err != nil {
		return err
	}
	if !ok {
		return errs.New(errs.CodeNoContent, "invalid parent id", nil)
	}
	return nil
}

func (s *imageFolderService) ensureNoDuplicateName(ctx context.Context, parentID *uuid.UUID, name string, excludeID *uuid.UUID) error {
	dup, err := s.imageFolderRepo.ExistsBySameNameInParent(ctx, s.db.Conn(), parentID, name, excludeID)
	if err != nil {
		return err
	}
	if dup {
		return errs.New(errs.CodeDuplicate, "duplicate name", nil)
	}
	return nil
}

func (s *imageFolderService) Create(ctx context.Context, name string, parentID *uuid.UUID) (*entity.ImageFolder, error) {
	name = normalizeName(name)
	if name == "" {
		return nil, errs.New(errs.CodeInvalidParam, "invalid name", nil)
	}

	if err := s.ensureParentExists(ctx, parentID); err != nil {
		return nil, err
	}
	if err := s.ensureNoDuplicateName(ctx, parentID, name, nil); err != nil {
		return nil, err
	}

	folder := &entity.ImageFolder{
		ID:       uuid.New(),
		ParentID: parentID,
		Name:     name,
	}
	if err := s.imageFolderRepo.Create(ctx, s.db.Conn(), folder); err != nil {
		return nil, err
	}
	return folder, nil
}

func (s *imageFolderService) Get(ctx context.Context, id uuid.UUID) (*entity.ImageFolder, error) {
	return s.ensureFolderExists(ctx, id)
}

func (s *imageFolderService) List(ctx context.Context, parentID *uuid.UUID, limit, offset int) ([]entity.ImageFolder, error) {
	if parentID != nil {
		if err := s.ensureParentExists(ctx, parentID); err != nil {
			return nil, err
		}
	}

	return s.imageFolderRepo.ListByParent(ctx, s.db.Conn(), parentID, limit, offset)
}

func (s *imageFolderService) Rename(ctx context.Context, id uuid.UUID, newName string) error {
	newName = normalizeName(newName)
	if newName == "" {
		return errs.New(errs.CodeInvalidParam, "invalid name", nil)
	}

	folder, err := s.ensureFolderExists(ctx, id)
	if err != nil {
		return err
	}

	if err := s.ensureNoDuplicateName(ctx, folder.ParentID, newName, &id); err != nil {
		return err
	}

	return s.imageFolderRepo.Rename(ctx, s.db.Conn(), id, newName)
}

func (s *imageFolderService) Move(ctx context.Context, id uuid.UUID, targetParentID *uuid.UUID) error {
	folder, err := s.ensureFolderExists(ctx, id)
	if err != nil {
		return err
	}

	// 允许移动到根目录，则无需校验父目录存在
	if err := s.ensureParentExists(ctx, targetParentID); err != nil {
		return err
	}

	// 防止把自己移动到自己下面
	if targetParentID != nil && *targetParentID == id {
		return errs.New(errs.CodeInvalidParam, "invalid target parent id", nil)
	}

	if err := s.ensureNoDuplicateName(ctx, targetParentID, folder.Name, &id); err != nil {
		return err
	}

	return s.imageFolderRepo.Move(ctx, s.db.Conn(), id, *targetParentID)
}

func (s *imageFolderService) Delete(ctx context.Context, id uuid.UUID) error {
	conn := s.db.Conn()

	_, err := s.ensureFolderExists(ctx, id)
	if err != nil {
		return err
	}
	return s.imageFolderRepo.SoftDelete(ctx, conn, id)
}

func NewImageFolderService(db database.DB, imageFolderRepo repository.IImageFolderRepo) IImageFolderService {
	return &imageFolderService{
		db:              db,
		imageFolderRepo: imageFolderRepo,
	}
}
