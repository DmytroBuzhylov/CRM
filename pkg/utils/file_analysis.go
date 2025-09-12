package utils

import (
	"Test/pkg/storage"
	"errors"
	"mime/multipart"
	"net/http"
)

func FileVerifier(file multipart.File, fileHeader *multipart.FileHeader) (*storage.StorageFileOptions, error) {
	if fileHeader.Size > 1024*1024*2 {
		return nil, errors.New("file is too large. Max size is 2MB")
	}
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return nil, errors.New("failed to read file buffer")
	}
	file.Seek(0, 0)

	contentType := http.DetectContentType(buffer)
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		return nil, errors.New("unsupported file type. Only JPEG, PNG, and GIF are allowed")
	}

	fileOpts := &storage.StorageFileOptions{
		File:        file,
		FileSize:    fileHeader.Size,
		ContentType: contentType,
	}

	return fileOpts, nil
}
