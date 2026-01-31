package service

import (
	"context"
	"io"

	"blog-server/storage"
)

type IImageService interface {
	Upload(ctx context.Context, key string, body io.Reader, contentType string) error
	Download(ctx context.Context, key string) (io.ReadCloser, string, error)
	Delete(ctx context.Context, key string) error
}

type imageService struct {
	store  storage.Storage
	bucket string
}

func (i *imageService) Delete(ctx context.Context, key string) error {
	return i.store.Delete(ctx, i.bucket, key)
}

func (i *imageService) Download(ctx context.Context, key string) (io.ReadCloser, string, error) {
	return i.store.Download(ctx, i.bucket, key)
}

func (i *imageService) Upload(ctx context.Context, key string, body io.Reader, contentType string) error {
	return i.store.Upload(ctx, i.bucket, key, body, contentType)
}

func NewImageService(storage storage.Storage) IImageService {
	return &imageService{store: storage, bucket: "image"}
}
