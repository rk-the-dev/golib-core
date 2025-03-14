package fileuploader

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// FileUploader struct handles file uploads
type FileUploader struct {
	UploadDir    string
	MaxSize      int64    // Max file size in bytes
	AllowedTypes []string // List of allowed MIME types
}

// New creates a new FileUploader instance
func New(uploadDir string, maxSize int64, allowedTypes []string) *FileUploader {
	return &FileUploader{
		UploadDir:    uploadDir,
		MaxSize:      maxSize,
		AllowedTypes: allowedTypes,
	}
}

// UploadFile handles the file upload
func (fu *FileUploader) UploadFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	if header.Size > fu.MaxSize {
		return "", fmt.Errorf("file size exceeds the limit of %d bytes", fu.MaxSize)
	}
	filePath := filepath.Join(fu.UploadDir, header.Filename)
	// Check if file type is allowed
	if !fu.isAllowedType(header) {
		return "", fmt.Errorf("file type not allowed")
	}
	outFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, file)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

// isAllowedType checks if the file type is allowed
func (fu *FileUploader) isAllowedType(header *multipart.FileHeader) bool {
	for _, t := range fu.AllowedTypes {
		if t == header.Header.Get("Content-Type") {
			return true
		}
	}
	return false
}
