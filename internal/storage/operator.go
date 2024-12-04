package storage

import (
	"io"
	"mime/multipart"
)

type FileOperator interface {
	SaveFile(file multipart.File, filename string) (string, error)
	DeleteFile(filepath string) error
	GetFile(filepath string) (io.ReadCloser, error)
	GenerateStoragePath(originalName string) string
	ValidateFile(file *multipart.FileHeader) error
}
