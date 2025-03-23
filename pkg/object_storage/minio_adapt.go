package objStorage

import (
	"blog/config"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioClient(cfg config.MinIO) (*minio.Client, error) {
	finalEndPoint := fmt.Sprintf("%s:%d", cfg.Endpoint, cfg.Port)

	minioClient, err := minio.New(finalEndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})

	if err != nil {
		return nil, err
	}

	cancel,err := minioClient.HealthCheck(time.Minute)
	if err != nil {
		return nil, err
	}

	defer cancel()

	return minioClient, nil
}
