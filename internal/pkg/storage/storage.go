package storage

import (
	"context"
	"io"
	"time"
)

type Storage interface {
	GeneratePresignedURL(ctx context.Context, fileID string, expiry time.Duration) (string, error)
	UploadFile(ctx context.Context, fileID string, reader io.Reader, fileSize int64, contentType string) error
	DeleteFile(ctx context.Context, fileID string) error
	//CheckFileExists(ctx context.Context, fileID string) (bool, string)
	GetFileReader(ctx context.Context, fileID string) (io.Reader, error)
}

type StorageFileOptions struct {
	File        io.Reader
	FileSize    int64
	ContentType string
}
