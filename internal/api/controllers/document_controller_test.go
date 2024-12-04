package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zhaoyikaiii/docmind/internal/models"
	"github.com/Zhaoyikaiii/docmind/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDocumentService 模拟文档服务
type MockDocumentService struct {
	mock.Mock
}

func (m *MockDocumentService) CreateDocument(ctx context.Context, doc *models.Document) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

func (m *MockDocumentService) UpdateDocument(ctx context.Context, doc *models.Document) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

func (m *MockDocumentService) DeleteDocument(ctx context.Context, id uint, userID uint) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *MockDocumentService) GetDocument(ctx context.Context, id uint) (*models.Document, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Document), args.Error(1)
}

func (m *MockDocumentService) ListDocuments(ctx context.Context, params repository.DocumentListParams) ([]models.Document, int64, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]models.Document), args.Get(1).(int64), args.Error(2)
}

func (m *MockDocumentService) CreateVersion(ctx context.Context, docID uint, userID uint) error {
	args := m.Called(ctx, docID, userID)
	return args.Error(0)
}

func (m *MockDocumentService) GetVersions(ctx context.Context, docID uint) ([]models.DocumentVersion, error) {
	args := m.Called(ctx, docID)
	return args.Get(0).([]models.DocumentVersion), args.Error(1)
}

func (m *MockDocumentService) ManageTags(ctx context.Context, docID uint, addTags []uint, removeTags []uint) error {
	args := m.Called(ctx, docID, addTags, removeTags)
	return args.Error(0)
}

func setupTest() (*gin.Engine, *MockDocumentService) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockDocumentService)
	controller := NewDocumentController(mockService)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	})

	docs := r.Group("/documents")
	{
		docs.POST("", controller.CreateDocument)
		docs.PUT("/:id", controller.UpdateDocument)
		docs.DELETE("/:id", controller.DeleteDocument)
		docs.GET("/:id", controller.GetDocument)
		docs.GET("", controller.ListDocuments)
		docs.POST("/:id/versions", controller.CreateVersion)
		docs.GET("/:id/versions", controller.GetVersions)
		docs.POST("/:id/tags", controller.ManageTags)
	}

	return r, mockService
}

func TestCreateDocument(t *testing.T) {
	r, mockService := setupTest()

	tests := []struct {
		name         string
		document     models.Document
		setupMock    func()
		expectedCode int
		description  string
	}{
		{
			name: "Success_CreateNewDocument",
			document: models.Document{
				Title:   "Test Document",
				Content: "Test Content",
			},
			setupMock: func() {
				mockService.On("CreateDocument", mock.Anything, mock.AnythingOfType("*models.Document")).
					Return(nil)
			},
			expectedCode: http.StatusCreated,
			description:  "应该成功创建新文档并返回201状态码",
		},
		{
			name: "Error_DuplicateTitle",
			document: models.Document{
				Title:   "Existing Document",
				Content: "Content",
			},
			setupMock: func() {
				mockService.On("CreateDocument", mock.Anything, mock.AnythingOfType("*models.Document")).
					Return(fmt.Errorf("document with this title already exists"))
			},
			expectedCode: http.StatusBadRequest,
			description:  "当尝试创建标题重复的文档时应该返回400状态码",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 清理之前的 mock
			mockService.ExpectedCalls = nil
			mockService.Calls = nil
			
			tt.setupMock()

			jsonDoc, _ := json.Marshal(tt.document)
			req, _ := http.NewRequest(http.MethodPost, "/documents", bytes.NewBuffer(jsonDoc))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code, 
				"期望状态码 %d，但得到 %d\n响应体: %s", 
				tt.expectedCode, w.Code, w.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetDocument(t *testing.T) {
	r, mockService := setupTest()

	tests := []struct {
		name         string
		documentID   string
		setupMock    func()
		expectedCode int
	}{
		{
			name:       "Valid document retrieval",
			documentID: "1",
			setupMock: func() {
				mockService.On("GetDocument", mock.Anything, uint(1)).
					Return(&models.Document{ID: 1, Title: "Test"}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:       "Document not found",
			documentID: "999",
			setupMock: func() {
				mockService.On("GetDocument", mock.Anything, uint(999)).
					Return(nil, fmt.Errorf("not found"))
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Invalid document ID",
			documentID:   "invalid",
			setupMock:    func() {},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req, _ := http.NewRequest(http.MethodGet, "/documents/"+tt.documentID, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestListDocuments(t *testing.T) {
	r, mockService := setupTest()

	tests := []struct {
		name         string
		query        string
		setupMock    func()
		expectedCode int
	}{
		{
			name:  "List documents with pagination",
			query: "?page=1&page_size=10",
			setupMock: func() {
				mockService.On("ListDocuments", mock.Anything, mock.AnythingOfType("repository.DocumentListParams")).
					Return([]models.Document{{ID: 1}}, int64(1), nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:  "List documents with search",
			query: "?search=test",
			setupMock: func() {
				mockService.On("ListDocuments", mock.Anything, mock.AnythingOfType("repository.DocumentListParams")).
					Return([]models.Document{}, int64(0), nil)
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req, _ := http.NewRequest(http.MethodGet, "/documents"+tt.query, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestManageTags(t *testing.T) {
	r, mockService := setupTest()

	tests := []struct {
		name         string
		documentID   string
		request      map[string][]uint
		setupMock    func()
		expectedCode int
	}{
		{
			name:       "Add and remove tags",
			documentID: "1",
			request: map[string][]uint{
				"add_tags":    {1, 2},
				"remove_tags": {3, 4},
			},
			setupMock: func() {
				mockService.On("ManageTags", mock.Anything, uint(1), []uint{1, 2}, []uint{3, 4}).
					Return(nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:       "Invalid document ID",
			documentID: "invalid",
			request: map[string][]uint{
				"add_tags": {1},
			},
			setupMock:    func() {},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			jsonReq, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest(http.MethodPost, "/documents/"+tt.documentID+"/tags", bytes.NewBuffer(jsonReq))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}
