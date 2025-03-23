package objStorage

import (
	"bytes"
	"context"
	"time"

	"github.com/minio/minio-go/v7"
)

type StorageRepository interface {
	UploadFile(ctx context.Context, bucketName string, object Object) error
	// DownloadFile(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, bucketName, objectName string) error
	GetObjectURL(ctx context.Context, bucketName, objectName string) (string, error)
}

const ExpirationTime = time.Second*60*60

type minioStorageRepo struct {
	client *minio.Client
}

func NewMinioStorageRepo(client *minio.Client) StorageRepository {
	return &minioStorageRepo{
		client: client,
	}
}

func (m *minioStorageRepo) UploadFile(ctx context.Context, bucketName string, object Object) error {
	file := bytes.NewReader(object.Data)
	objectPath := "profile-images/" + object.Name

	_, err := m.client.PutObject(ctx, bucketName, objectPath, file, object.Size, minio.PutObjectOptions{ContentType: object.ContentType})
	if err != nil {
		return err
	}

	return nil
}


func (m *minioStorageRepo) GetObjectURL(ctx context.Context, bucketName, objectName string) (string, error) {
	objectPath := "profile-images/" + objectName

	// Check if the object exists
	_, err := m.client.StatObject(ctx, bucketName, objectPath, minio.StatObjectOptions{})
	if err != nil {
		// If StatObject returns an error, it's likely the object doesn't exist.
		// You can check the error type to be more precise (minio.ErrorResponse).
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return "", ErrObjectNotFound 
		}
		return "", err 
	}

	url, err := m.client.PresignedGetObject(ctx, bucketName, objectPath, ExpirationTime, nil)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}

func (m *minioStorageRepo) DeleteFile(ctx context.Context, bucketName, objectName string) error {
	objectPath := "profile-images/" + objectName

	err := m.client.RemoveObject(ctx, bucketName, objectPath, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
