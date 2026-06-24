package adapter

import (
	"bytes"
	"context"
	"deep-map-server/internal/domain"
	objectstorage "deep-map-server/internal/infrastructure/object_storage"
)

// ObjectOSSAdapter 将 AWS S3 SDK 客户端适配为 domain.ObjectStorage 端口。
type ObjectOSSAdapter struct {
	client *objectstorage.Client
}

func NewObjectOSSAdapter(client *objectstorage.Client) *ObjectOSSAdapter {
	return &ObjectOSSAdapter{client: client}
}

var _ domain.ObjectStorage = (*ObjectOSSAdapter)(nil)

func (a *ObjectOSSAdapter) Upload(ctx context.Context, fileKey string, file []byte, contentType string) error {
	return a.client.Upload(ctx, fileKey, bytes.NewReader(file), contentType)
}

func (a *ObjectOSSAdapter) Download(ctx context.Context, fileKey string) ([]byte, error) {
	data, _, err := a.client.GetObject(ctx, fileKey)
	return data, err
}

func (a *ObjectOSSAdapter) Delete(ctx context.Context, fileKey string) error {
	return a.client.DeleteObject(ctx, fileKey)
}

func (a *ObjectOSSAdapter) PresignGetObjectURL(ctx context.Context, objectKey string) (string, error) {
	return a.client.PresignGetObjectURL(ctx, objectKey)
}
