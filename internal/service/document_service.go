package service

import (
	"context"
	"errors"
	"github.com/Zhaoyikaiii/docmind/internal/models"
	"github.com/Zhaoyikaiii/docmind/internal/repository"
)

type DocumentService interface {
	CreateDocument(ctx context.Context, doc *models.Document) error
	UpdateDocument(ctx context.Context, doc *models.Document) error
	DeleteDocument(ctx context.Context, id uint, userID uint) error
	GetDocument(ctx context.Context, id uint) (*models.Document, error)
	ListDocuments(ctx context.Context, params repository.DocumentListParams) ([]models.Document, int64, error)
	CreateVersion(ctx context.Context, docID uint, userID uint) error
	GetVersions(ctx context.Context, docID uint) ([]models.DocumentVersion, error)
	ManageTags(ctx context.Context, docID uint, addTags []uint, removeTags []uint) error
}

type documentService struct {
	repo repository.DocumentRepository
}

func NewDocumentService(repo repository.DocumentRepository) DocumentService {
	return &documentService{repo: repo}
}

func (s *documentService) CreateDocument(ctx context.Context, doc *models.Document) error {
	return s.repo.Create(ctx, doc)
}

func (s *documentService) UpdateDocument(ctx context.Context, doc *models.Document) error {
	existing, err := s.repo.GetByID(ctx, doc.ID)
	if err != nil {
		return err
	}

	if existing.CreatorID != doc.CreatorID {
		return errors.New("unauthorized to update this document")
	}

	return s.repo.Update(ctx, doc)
}

func (s *documentService) DeleteDocument(ctx context.Context, id uint, userID uint) error {
	doc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if doc.CreatorID != userID {
		return errors.New("unauthorized to delete this document")
	}

	return s.repo.Delete(ctx, id)
}

func (s *documentService) GetDocument(ctx context.Context, id uint) (*models.Document, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *documentService) ListDocuments(ctx context.Context, params repository.DocumentListParams) ([]models.Document, int64, error) {
	return s.repo.List(ctx, params)
}

func (s *documentService) CreateVersion(ctx context.Context, docID uint, userID uint) error {
	doc, err := s.repo.GetByID(ctx, docID)
	if err != nil {
		return err
	}

	version := &models.DocumentVersion{
		DocumentID: docID,
		Version:    doc.Version,
		Title:      doc.Title,
		Content:    doc.Content,
		CreatedBy:  userID,
	}

	return s.repo.CreateVersion(ctx, version)
}

func (s *documentService) GetVersions(ctx context.Context, docID uint) ([]models.DocumentVersion, error) {
	return s.repo.GetVersions(ctx, docID)
}

func (s *documentService) ManageTags(ctx context.Context, docID uint, addTags []uint, removeTags []uint) error {
	if len(addTags) > 0 {
		if err := s.repo.AddTags(ctx, docID, addTags); err != nil {
			return err
		}
	}

	if len(removeTags) > 0 {
		if err := s.repo.RemoveTags(ctx, docID, removeTags); err != nil {
			return err
		}
	}

	return nil
}
