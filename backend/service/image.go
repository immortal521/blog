package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"blog-server/database"
	"blog-server/entity"
	"blog-server/errs"
	"blog-server/repository"
	"blog-server/storage"

	"github.com/google/uuid"
)

type IImageService interface {
	Upload(ctx context.Context, folderID *uuid.UUID, filename string, body io.Reader, contentType string, size int64) (*entity.Image, error)
	Download(ctx context.Context, key string) (io.ReadCloser, string, error)
	Delete(ctx context.Context, key string) error
}

type imageService struct {
	db         database.DB
	store      storage.Storage
	bucket     string
	folderRepo repository.IImageFolderRepo
	imageRepo  repository.IImageRepo
}

func (i *imageService) Delete(ctx context.Context, key string) error {
	return i.store.Delete(ctx, i.bucket, key)
}

func (i *imageService) Download(ctx context.Context, key string) (io.ReadCloser, string, error) {
	return i.store.Download(ctx, i.bucket, key)
}

func (i *imageService) Upload(ctx context.Context, folderID *uuid.UUID, filename string, body io.Reader, contentType string, size int64) (*entity.Image, error) {
	filename = strings.TrimSpace(filename)
	contentType = strings.TrimSpace(contentType)

	if filename == "" {
		return nil, errs.New(errs.CodeInvalidParam, "filename is empty", nil)
	}
	if body == nil {
		return nil, errs.New(errs.CodeInvalidParam, "invalid body", nil)
	}
	if size <= 0 {
		return nil, errs.New(errs.CodeInvalidParam, "invalid size", nil)
	}

	if !strings.HasPrefix(strings.ToLower(contentType), "image/") {
		return nil, errs.New(errs.CodeInvalidParam, "unsupported content type", nil)
	}

	if folderID != nil {
		ok, err := i.folderRepo.Exists(ctx, i.db.Conn(), *folderID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errs.New(errs.CodeNoContent, "no such folder", nil)
		}
	}

	tmp, err := os.CreateTemp("", "img-upload-*")
	if err != nil {
		return nil, err
	}
	tmpPath := tmp.Name()
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
	}()

	written, err := io.Copy(tmp, body)
	if err != nil {
		return nil, err
	}

	if _, err := tmp.Seek(0, 0); err != nil {
		return nil, err
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, tmp); err != nil {
		return nil, err
	}

	hashHex := hex.EncodeToString(hasher.Sum(nil))

	ext := strings.ToLower(filepath.Ext(filename))
	finalKey := fmt.Sprintf("%s/%s%s", hashHex[:2], hashHex, ext)

	exist, err := i.imageRepo.GetBySha256(ctx, i.db.Conn(), hashHex)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		now := time.Now()
		img := &entity.Image{
			ID:         uuid.New(),
			FolderID:   folderID,
			StorageKey: exist.StorageKey,
			OriginName: filename,
			Mime:       contentType,
			Size:       written,
			Sha256:     &hashHex,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		if err := i.imageRepo.Create(ctx, i.db.Conn(), img); err != nil {
			return nil, err
		}
		return img, nil
	}

	ok, err := i.store.Exists(ctx, i.bucket, finalKey)
	if err != nil {
		return nil, err
	}
	if !ok {
		if _, err := tmp.Seek(0, 0); err != nil {
			return nil, err
		}
		if err := i.store.Upload(ctx, i.bucket, finalKey, tmp, contentType); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	img := &entity.Image{
		ID:         uuid.New(),
		FolderID:   folderID,
		StorageKey: finalKey,
		OriginName: filename,
		Mime:       contentType,
		Size:       written,
		Sha256:     &hashHex,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := i.imageRepo.Create(ctx, i.db.Conn(), img); err != nil {
		return nil, err
	}
	return img, nil
}

func NewImageService(db database.DB, storage storage.Storage, folderRepo repository.IImageFolderRepo, imageRepo repository.IImageRepo) IImageService {
	return &imageService{db: db, store: storage, bucket: "images", folderRepo: folderRepo, imageRepo: imageRepo}
}
