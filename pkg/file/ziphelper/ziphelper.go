package ziphelper

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// ZipHelper defines the interface for file compression operations
type ZipHelper interface {
	Zip(source, destination string) error
	Unzip(source, destination string) error
	Close() error
}

// zipHelper implements ZipHelper
type zipHelper struct{}

var (
	instance *zipHelper
	once     sync.Once
)

// NewZipHelper initializes and returns a ZipHelper instance
func NewZipHelper() ZipHelper {
	once.Do(func() {
		instance = &zipHelper{}
	})
	return instance
}

// Zip compresses a file or directory into a .zip file
func (z *zipHelper) Zip(source, destination string) error {
	outFile, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	// Walk through the source and add files to the archive
	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create zip header
		relPath, err := filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath

		if info.IsDir() {
			header.Name += "/"
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to zip files: %w", err)
	}

	fmt.Println("✅ Zipped successfully:", destination)
	return nil
}

// Unzip extracts a .zip archive to the specified destination
func (z *zipHelper) Unzip(source, destination string) error {
	zipReader, err := zip.OpenReader(source)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		filePath := filepath.Join(destination, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		if err != nil {
			return err
		}

		outFile, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		srcFile, err := file.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()

		_, err = io.Copy(outFile, srcFile)
		if err != nil {
			return err
		}
	}

	fmt.Println("✅ Unzipped successfully:", destination)
	return nil
}

// Close is a placeholder to maintain interface consistency
func (z *zipHelper) Close() error {
	return nil
}
