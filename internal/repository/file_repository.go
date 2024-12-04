package repository

import (
	"context"

	"github.com/Zhaoyikaiii/docmind/internal/models"
	"gorm.io/gorm"
)

type FileRepository interface {
	Create(ctx context.Context, file *models.File) error
	Update(ctx context.Context, file *models.File) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*models.File, error)
	List(ctx context.Context, params FileListParams) ([]models.File, int64, error)
}

type FileListParams struct {
	UploaderID  *uint
	DocumentID  *uint
	ContentType *string
	Search      string
	Page        int
	PageSize    int
}

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) Create(ctx context.Context, file *models.File) error {
	return r.db.WithContext(ctx).Create(file).Error
}

func (r *fileRepository) Update(ctx context.Context, file *models.File) error {
	return r.db.WithContext(ctx).Save(file).Error
}

func (r *fileRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.File{}, id).Error
}

func (r *fileRepository) GetByID(ctx context.Context, id uint) (*models.File, error) {
	var file models.File
	err := r.db.WithContext(ctx).First(&file, id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) List(ctx context.Context, params FileListParams) ([]models.File, int64, error) {
	var files []models.File
	var total int64

	query := r.db.WithContext(ctx).Model(&models.File{})

	if params.UploaderID != nil {
		query = query.Where("uploader_id = ?", *params.UploaderID)
	}

	if params.DocumentID != nil {
		query = query.Where("document_id = ?", *params.DocumentID)
	}

	if params.ContentType != nil {
		query = query.Where("content_type = ?", *params.ContentType)
	}

	if params.Search != "" {
		query = query.Where("original_name LIKE ?", "%"+params.Search+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((params.Page - 1) * params.PageSize).
		Limit(params.PageSize).
		Find(&files).Error

	return files, total, err
} 