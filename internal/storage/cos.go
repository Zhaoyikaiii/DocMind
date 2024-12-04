package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type COSFileOperator struct {
	client *cos.Client
	config COSConfig
}

type COSConfig struct {
	Region       string
	BucketName   string
	SecretID     string
	SecretKey    string
	BasePath     string
	MaxFileSize  int64
	AllowedTypes map[string]bool
}

func NewCOSFileOperator(config COSConfig) (FileOperator, error) {
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com",
		config.BucketName, config.Region))
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.SecretID,
			SecretKey: config.SecretKey,
		},
	})

	return &COSFileOperator{
		client: client,
		config: config,
	}, nil
}

func (c *COSFileOperator) SaveFile(file multipart.File, filename string) (string, error) {
	objectKey := c.GenerateStoragePath(filename)

	_, err := c.client.Object.Put(context.Background(), objectKey, file, nil)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to COS: %w", err)
	}

	return c.client.Object.GetObjectURL(objectKey).String(), nil
}

func (c *COSFileOperator) DeleteFile(filepath string) error {
	objectKey := c.getObjectKeyFromPath(filepath)

	_, err := c.client.Object.Delete(context.Background(), objectKey)
	if err != nil {
		return fmt.Errorf("failed to delete file from COS: %w", err)
	}

	return nil
}

func (c *COSFileOperator) GetFile(filepath string) (io.ReadCloser, error) {
	objectKey := c.getObjectKeyFromPath(filepath)

	resp, err := c.client.Object.Get(context.Background(), objectKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get file from COS: %w", err)
	}

	return resp.Body, nil
}

func (c *COSFileOperator) GenerateStoragePath(originalName string) string {
	ext := filepath.Ext(originalName)
	uniqueName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	return filepath.Join(
		c.config.BasePath,
		time.Now().Format("2006/01/02"),
		uniqueName,
	)
}

func (c *COSFileOperator) ValidateFile(file *multipart.FileHeader) error {
	// 验证文件大小
	if file.Size > c.config.MaxFileSize {
		return fmt.Errorf("file too large, maximum size is %d bytes", c.config.MaxFileSize)
	}

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !c.config.AllowedTypes[ext] {
		return fmt.Errorf("file type %s not allowed", ext)
	}

	return nil
}

// 从完整URL中提取对象键
func (c *COSFileOperator) getObjectKeyFromPath(path string) string {
	prefix := fmt.Sprintf("https://%s.cos.%s.myqcloud.com/",
		c.config.BucketName, c.config.Region)
	return strings.TrimPrefix(path, prefix)
}
