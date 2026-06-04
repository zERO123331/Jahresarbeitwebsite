package cdn

import (
	"Jahresarbeitwebsite/internal/validator"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

const (
	TypeJpeg = "image/jpeg"
	TypePng  = "image/png"
	TypeGif  = "image/gif"
	TypeWebp = "image/webp"
)

type CDN struct {
	client *minio.Client
	bucket string
}

func IsAllowedImageType(contentType string) bool {
	switch contentType {
	case TypeJpeg, TypePng, TypeGif, TypeWebp:
		return true
	default:
		return false
	}
}

func New(client *minio.Client, bucket string) *CDN {
	return &CDN{
		client: client,
		bucket: bucket,
	}
}

func IsAllowedContentType(contentType string) bool {
	return validator.PermittedValue(contentType, TypeJpeg, TypePng, TypeGif, TypeWebp)
}

func (c *CDN) Upload(ctx context.Context, file io.Reader, size int64, objectKey, contentType string) (string, error) {
	if file == nil {
		return "", fmt.Errorf("file is nil")
	}

	_, err := c.client.PutObject(ctx, c.bucket, objectKey, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("/%s", objectKey), nil
}

func (c *CDN) Delete(ctx context.Context, objectKey string) error {
	err := c.client.RemoveObject(ctx, c.bucket, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
