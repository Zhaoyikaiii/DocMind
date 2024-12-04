package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type QiniuFileOperator struct {
	mac    *qbox.Mac
	bucket string
	config QiniuConfig
}

type QiniuConfig struct {
	AccessKey    string
	SecretKey    string
	BucketName   string
	Domain       string
	BasePath     string
	MaxFileSize  int64
	AllowedTypes map[string]bool
}

func NewQiniuFileOperator(config QiniuConfig) FileOperator {
	return &QiniuFileOperator{
		mac:    qbox.NewMac(config.AccessKey, config.SecretKey),
		bucket: config.BucketName,
		config: config,
	}
}

func (q *QiniuFileOperator) SaveFile(file multipart.File, filename string) (string, error) {
	objectKey := q.GenerateStoragePath(filename)

	putPolicy := storage.PutPolicy{
		Scope: q.bucket,
	}
	upToken := putPolicy.UploadToken(q.mac)

	cfg := storage.Config{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	err := formUploader.Put(context.Background(), &ret, upToken, objectKey, file, -1, nil)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to Qiniu: %w", err)
	}

	return fmt.Sprintf("https://%s/%s", q.config.Domain, ret.Key), nil
}

func (q *QiniuFileOperator) DeleteFile(filepath string) error {
	objectKey := q.getObjectKeyFromPath(filepath)

	bucketManager := storage.NewBucketManager(q.mac, nil)
	err := bucketManager.Delete(q.bucket, objectKey)
	if err != nil {
		return fmt.Errorf("failed to delete file from Qiniu: %w", err)
	}

	return nil
}

func (q *QiniuFileOperator) GetFile(filepath string) (io.ReadCloser, error) {
	objectKey := q.getObjectKeyFromPath(filepath)

	// 生成私有下载链接
	deadline := time.Now().Add(time.Hour).Unix()
	privateAccessURL := storage.MakePrivateURL(q.mac, q.config.Domain, objectKey, deadline)

	// 获取文件内容
	resp, err := http.Get(privateAccessURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get file from Qiniu: %w", err)
	}

	return resp.Body, nil
}

func (q *QiniuFileOperator) GenerateStoragePath(originalName string) string {
	ext := filepath.Ext(originalName)
	uniqueName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	return filepath.Join(
		q.config.BasePath,
		time.Now().Format("2006/01/02"),
		uniqueName,
	)
}

func (q *QiniuFileOperator) ValidateFile(file *multipart.FileHeader) error {
	if file.Size > q.config.MaxFileSize {
		return fmt.Errorf("file too large, maximum size is %d bytes", q.config.MaxFileSize)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !q.config.AllowedTypes[ext] {
		return fmt.Errorf("file type %s not allowed", ext)
	}

	return nil
}

func (q *QiniuFileOperator) getObjectKeyFromPath(path string) string {
	prefix := fmt.Sprintf("https://%s/", q.config.Domain)
	return strings.TrimPrefix(path, prefix)
}
