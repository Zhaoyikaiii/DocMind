package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Zhaoyikaiii/docmind/internal/models"
	"github.com/google/uuid"
)

type LocalFileOperator struct {
	config LocalConfig
}

// LocalConfig 本地存储配置
type LocalConfig struct {
	UploadDir    string
	MaxFileSize  int64
	AllowedTypes map[string]bool
}

func NewLocalFileOperator(config LocalConfig) models.FileOperator {
	return &LocalFileOperator{
		config: config,
	}
}

func (l *LocalFileOperator) SaveFile(file multipart.File, filename string) (string, error) {
	storagePath := l.GenerateStoragePath(filename)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(storagePath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// 创建文件
	dst, err := os.Create(storagePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err = io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return storagePath, nil
}

func (l *LocalFileOperator) DeleteFile(filepath string) error {
	if err := os.Remove(filepath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (l *LocalFileOperator) GetFile(filepath string) (io.ReadCloser, error) {
	return os.Open(filepath)
}

func (l *LocalFileOperator) GenerateStoragePath(originalName string) string {
	ext := filepath.Ext(originalName)
	uniqueName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	return filepath.Join(l.config.UploadDir,
		time.Now().Format("2006/01/02"),
		uniqueName)
}

func (l *LocalFileOperator) ValidateFile(file *multipart.FileHeader) error {
	// 验证文件大小
	if file.Size > l.config.MaxFileSize {
		return fmt.Errorf("file too large, maximum size is %d bytes", l.config.MaxFileSize)
	}

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !l.config.AllowedTypes[ext] {
		return fmt.Errorf("file type %s not allowed", ext)
	}

	return nil
}
