package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
)

type OSSFileOperator struct {
	client *oss.Client
	bucket *oss.Bucket
	config OSSConfig
}

type OSSConfig struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
	BasePath        string
	MaxFileSize     int64
	AllowedTypes    map[string]bool
}

func NewOSSFileOperator(config OSSConfig) (FileOperator, error) {
	// 创建OSS客户端
	client, err := oss.New(config.Endpoint, config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create OSS client: %w", err)
	}

	// 获取存储空间
	bucket, err := client.Bucket(config.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket: %w", err)
	}

	return &OSSFileOperator{
		client: client,
		bucket: bucket,
		config: config,
	}, nil
}

func (o *OSSFileOperator) SaveFile(file multipart.File, filename string) (string, error) {
	objectKey := o.GenerateStoragePath(filename)

	// 上传文件到OSS
	err := o.bucket.PutObject(objectKey, file)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to OSS: %w", err)
	}

	// 返回文件的完整URL或OSS路径
	return fmt.Sprintf("https://%s.%s/%s",
		o.config.BucketName,
		o.config.Endpoint,
		objectKey), nil
}

func (o *OSSFileOperator) DeleteFile(filepath string) error {
	// 从URL中提取objectKey
	objectKey := o.getObjectKeyFromPath(filepath)

	// 删除OSS对象
	err := o.bucket.DeleteObject(objectKey)
	if err != nil {
		return fmt.Errorf("failed to delete file from OSS: %w", err)
	}

	return nil
}

func (o *OSSFileOperator) GetFile(filepath string) (io.ReadCloser, error) {
	objectKey := o.getObjectKeyFromPath(filepath)

	// 获取文件流
	reader, err := o.bucket.GetObject(objectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get file from OSS: %w", err)
	}

	return reader, nil
}

func (o *OSSFileOperator) GenerateStoragePath(originalName string) string {
	ext := filepath.Ext(originalName)
	uniqueName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	return filepath.Join(
		o.config.BasePath,
		time.Now().Format("2006/01/02"),
		uniqueName,
	)
}

func (o *OSSFileOperator) ValidateFile(file *multipart.FileHeader) error {
	// 验证文件大小
	if file.Size > o.config.MaxFileSize {
		return fmt.Errorf("file too large, maximum size is %d bytes", o.config.MaxFileSize)
	}

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !o.config.AllowedTypes[ext] {
		return fmt.Errorf("file type %s not allowed", ext)
	}

	return nil
}

// 从完整URL中提取对象键
func (o *OSSFileOperator) getObjectKeyFromPath(path string) string {
	prefix := fmt.Sprintf("https://%s.%s/", o.config.BucketName, o.config.Endpoint)
	return strings.TrimPrefix(path, prefix)
}
