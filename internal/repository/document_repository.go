package repository

import (
	"context"
	"github.com/Zhaoyikaiii/docmind/internal/models"
	"gorm.io/gorm"
)

type DocumentRepository interface {
	Create(ctx context.Context, doc *models.Document) error
	Update(ctx context.Context, doc *models.Document) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*models.Document, error)
	List(ctx context.Context, params DocumentListParams) ([]models.Document, int64, error)
	ExistsByTitleAndCreator(ctx context.Context, title string, creatorID uint) (bool, error)
	CreateVersion(ctx context.Context, version *models.DocumentVersion) error
	GetVersions(ctx context.Context, documentID uint) ([]models.DocumentVersion, error)
	AddTags(ctx context.Context, docID uint, tagIDs []uint) error
	RemoveTags(ctx context.Context, docID uint, tagIDs []uint) error
}

type DocumentListParams struct {
	CreatorID *uint
	Status    *string
	Tags      []string
	Search    string
	Page      int
	PageSize  int
}

type documentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) DocumentRepository {
	return &documentRepository{db: db}
}

func (r *documentRepository) Create(ctx context.Context, doc *models.Document) error {
	return r.db.WithContext(ctx).Create(doc).Error
}

func (r *documentRepository) Update(ctx context.Context, doc *models.Document) error {
	return r.db.WithContext(ctx).Save(doc).Error
}

func (r *documentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Document{}, id).Error
}

func (r *documentRepository) GetByID(ctx context.Context, id uint) (*models.Document, error) {
	var doc models.Document
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Tags").
		First(&doc, id).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *documentRepository) List(ctx context.Context, params DocumentListParams) ([]models.Document, int64, error) {
	var docs []models.Document
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Document{})

	if params.CreatorID != nil {
		query = query.Where("creator_id = ?", *params.CreatorID)
	}

	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}

	if len(params.Tags) > 0 {
		query = query.Joins("JOIN document_tags ON documents.id = document_tags.document_id").
			Joins("JOIN tags ON document_tags.tag_id = tags.id").
			Where("tags.name IN ?", params.Tags)
	}

	if params.Search != "" {
		query = query.Where("title LIKE ? OR content LIKE ?",
			"%"+params.Search+"%", "%"+params.Search+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((params.Page - 1) * params.PageSize).
		Limit(params.PageSize).
		Preload("Creator").
		Preload("Tags").
		Find(&docs).Error

	return docs, total, err
}

func (r *documentRepository) CreateVersion(ctx context.Context, version *models.DocumentVersion) error {
	return r.db.WithContext(ctx).Create(version).Error
}

func (r *documentRepository) GetVersions(ctx context.Context, documentID uint) ([]models.DocumentVersion, error) {
	var versions []models.DocumentVersion
	err := r.db.WithContext(ctx).
		Where("document_id = ?", documentID).
		Order("version DESC").
		Find(&versions).Error
	return versions, err
}

func (r *documentRepository) AddTags(ctx context.Context, docID uint, tagIDs []uint) error {
	return r.db.WithContext(ctx).Exec(
		"INSERT INTO document_tags (document_id, tag_id) SELECT ? as document_id, unnest(?::int[]) as tag_id",
		docID, tagIDs).Error
}

func (r *documentRepository) RemoveTags(ctx context.Context, docID uint, tagIDs []uint) error {
	return r.db.WithContext(ctx).
		Table("document_tags").
		Where("document_id = ? AND tag_id IN ?", docID, tagIDs).
		Delete(nil).Error
}

func (r *documentRepository) ExistsByTitleAndCreator(ctx context.Context, title string, creatorID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Document{}).
		Where("title = ? AND creator_id = ?", title, creatorID).
		Count(&count).Error
	return count > 0, err
}
