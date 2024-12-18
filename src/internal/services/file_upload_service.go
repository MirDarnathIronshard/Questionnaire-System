package services

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/google/uuid"
)

type FileUploadService struct {
	config *models.FileUploadConfig
}

func NewFileUploadService(config *models.FileUploadConfig) *FileUploadService {
	return &FileUploadService{
		config: config,
	}
}

func (s *FileUploadService) Upload(file *multipart.FileHeader, allowedTypes []string) (string, error) {
	// Validate file size
	if file.Size > s.config.MaxFileSize {
		return "", fmt.Errorf("file size exceeds maximum limit of %d bytes", s.config.MaxFileSize)
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !s.isAllowedFileType(ext, allowedTypes) {
		return "", fmt.Errorf("file type %s is not allowed", ext)
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Create unique filename
	filename := uuid.New().String() + ext
	year := time.Now().Format("2006")
	month := time.Now().Format("01")
	uploadPath := filepath.Join(s.config.UploadDir, year, month)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", err
	}

	// Create destination file
	fullPath := filepath.Join(uploadPath, filename)
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy file contents
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	// Return relative path for storage
	relativePath := filepath.Join(year, month, filename)
	return relativePath, nil
}

func (s *FileUploadService) Delete(filePath string) error {
	if filePath == "" {
		return errors.New("file path is empty")
	}

	fullPath := filepath.Join(s.config.UploadDir, filePath)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return errors.New("file does not exist")
	}

	// Remove file
	return os.Remove(fullPath)
}

func (s *FileUploadService) isAllowedFileType(ext string, allowedTypes []string) bool {
	for _, allowed := range allowedTypes {
		if ext == allowed {
			return true
		}
	}
	return false
}

func (s *FileUploadService) GetURL(filePath string) string {
	return fmt.Sprintf("%s/%s", s.config.BaseURL, filePath)
}
