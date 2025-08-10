package storage

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"time"
)

type minioStorage struct {
	client *minio.Client
	bucket string
}

func NewMinioStorage(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*minioStorage, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	found, err := minioClient.BucketExists(context.Background(), bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check if bucket exists: %w", err)
	}

	if !found {
		err = minioClient.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &minioStorage{
		client: minioClient,
		bucket: bucket,
	}, nil
}

func (s *minioStorage) UploadFile(ctx context.Context, fileID string, reader io.Reader, fileSize int64, contentType string) error {
	_, err := s.client.PutObject(ctx, s.bucket, fileID, reader, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// GetFileReader обязательно закрывать Reader через defer
func (s *minioStorage) GetFileReader(ctx context.Context, fileID string) (io.Reader, error) {
	object, err := s.client.GetObject(ctx, s.bucket, fileID, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (s *minioStorage) GeneratePresignedURL(ctx context.Context, fileID string, expiry time.Duration) (string, error) {
	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucket, fileID, expiry, nil)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

func (s *minioStorage) DeleteFile(ctx context.Context, fileID string) error {
	err := s.client.RemoveObject(ctx, s.bucket, fileID, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
