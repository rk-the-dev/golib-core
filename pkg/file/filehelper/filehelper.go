package filehelper

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// FileHelper defines the interface for file operations
type FileHelper interface {
	ReadFile(filename string) (string, error)
	WriteFile(filename string, content string) error
	AppendToFile(filename string, content string) error
	DeleteFile(filename string) error
	CopyFile(src, dest string) error
	Close() error
}

// fileHelper implements FileHelper
type fileHelper struct{}

var (
	instance *fileHelper
	once     sync.Once
)

// NewFileHelper initializes and returns a file helper instance
func NewFileHelper() FileHelper {
	once.Do(func() {
		instance = &fileHelper{}
	})
	return instance
}

// ReadFile reads the entire content of a file
func (f *fileHelper) ReadFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(data), nil
}

// WriteFile writes content to a file (overwrites if exists)
func (f *fileHelper) WriteFile(filename string, content string) error {
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

// AppendToFile appends content to a file
func (f *fileHelper) AppendToFile(filename string, content string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for appending: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to append to file: %w", err)
	}
	return nil
}

// DeleteFile removes a file from the system
func (f *fileHelper) DeleteFile(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// CopyFile copies a file from src to dest
func (f *fileHelper) CopyFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

// Close is a placeholder to maintain interface consistency
func (f *fileHelper) Close() error {
	// No persistent resources to close, added for consistency
	return nil
}
