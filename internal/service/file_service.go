package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/Zhaoyikaiii/docmind/internal/models"
	"github.com/Zhaoyikaiii/docmind/internal/repository"
	"github.com/Zhaoyikaiii/docmind/internal/storage"
)

type FileService interface {
	CreateFile(ctx context.Context, file *models.File, uploadedFile multipart.File) error
	GetFile(ctx context.Context, id uint) (*models.File, error)
	DeleteFile(ctx context.Context, id uint, userID uint) error
	ListFiles(ctx context.Context, params models.FileListParams) ([]models.File, int64, error)
	UpdateFile(ctx context.Context, file *models.File) error
}

type fileService struct {
	repo         repository.FileRepository
	fileOperator storage.FileOperator
}

func NewFileService(repo repository.FileRepository, fileOperator storage.FileOperator) FileService {
	return &fileService{
		repo:         repo,
		fileOperator: fileOperator,
	}
}

func (s *fileService) CreateFile(ctx context.Context, file *models.File, uploadedFile multipart.File) error {
	storagePath, err := s.fileOperator.SaveFile(uploadedFile, file.OriginalName)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	file.Path = storagePath
	return s.repo.Create(ctx, file)
}

func (s *fileService) GetFile(ctx context.Context, id uint) (*models.File, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *fileService) DeleteFile(ctx context.Context, id uint, userID uint) error {
	file, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if file.UploaderID != userID {
		return errors.New("unauthorized to delete this file")
	}

	if err := s.fileOperator.DeleteFile(file.Path); err != nil {
		return fmt.Errorf("failed to delete physical file: %w", err)
	}

	return s.repo.Delete(ctx, id)
}

func (s *fileService) ListFiles(ctx context.Context, params models.FileListParams) ([]models.File, int64, error) {
	repoParams := repository.FileListParams{
		UploaderID:  params.UploaderID,
		DocumentID:  params.DocumentID,
		ContentType: params.ContentType,
		Search:      params.Search,
		Page:        params.Page,
		PageSize:    params.PageSize,
	}
	return s.repo.List(ctx, repoParams)
}

func (s *fileService) UpdateFile(ctx context.Context, file *models.File) error {
	existing, err := s.repo.GetByID(ctx, file.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing file: %w", err)
	}

	if existing.UploaderID != file.UploaderID {
		return errors.New("unauthorized to update this file")
	}

	existing.DocumentID = file.DocumentID

	return s.repo.Update(ctx, existing)
}
