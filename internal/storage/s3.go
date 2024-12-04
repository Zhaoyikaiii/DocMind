package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

type S3FileOperator struct {
	session  *session.Session
	uploader *s3manager.Uploader
	config   S3Config
}

type S3Config struct {
	Region          string
	BucketName      string
	AccessKeyID     string
	AccessKeySecret string
	BasePath        string
	MaxFileSize     int64
	AllowedTypes    map[string]bool
}

func NewS3FileOperator(config S3Config) (FileOperator, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(config.AccessKeyID, config.AccessKeySecret, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	return &S3FileOperator{
		session:  sess,
		uploader: s3manager.NewUploader(sess),
		config:   config,
	}, nil
}

func (s *S3FileOperator) SaveFile(file multipart.File, filename string) (string, error) {
	objectKey := s.GenerateStoragePath(filename)

	_, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.config.BucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
		s.config.BucketName,
		s.config.Region,
		objectKey), nil
}

func (s *S3FileOperator) DeleteFile(filepath string) error {
	objectKey := s.getObjectKeyFromPath(filepath)

	svc := s3.New(s.session)
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.config.BucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

func (s *S3FileOperator) GetFile(filepath string) (io.ReadCloser, error) {
	objectKey := s.getObjectKeyFromPath(filepath)

	svc := s3.New(s.session)
	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.config.BucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file from S3: %w", err)
	}

	return result.Body, nil
}

func (s *S3FileOperator) GenerateStoragePath(originalName string) string {
	ext := filepath.Ext(originalName)
	uniqueName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	return filepath.Join(
		s.config.BasePath,
		time.Now().Format("2006/01/02"),
		uniqueName,
	)
}

func (s *S3FileOperator) ValidateFile(file *multipart.FileHeader) error {
	if file.Size > s.config.MaxFileSize {
		return fmt.Errorf("file too large, maximum size is %d bytes", s.config.MaxFileSize)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !s.config.AllowedTypes[ext] {
		return fmt.Errorf("file type %s not allowed", ext)
	}

	return nil
}

func (s *S3FileOperator) getObjectKeyFromPath(path string) string {
	prefix := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/",
		s.config.BucketName,
		s.config.Region)
	return strings.TrimPrefix(path, prefix)
}
