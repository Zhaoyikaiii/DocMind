package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Zhaoyikaiii/docmind/internal/models"
	"github.com/Zhaoyikaiii/docmind/internal/repository"
	"github.com/Zhaoyikaiii/docmind/internal/service"
	"github.com/gin-gonic/gin"
)

type DocumentController struct {
	docService service.DocumentService
}

func NewDocumentController(docService service.DocumentService) *DocumentController {
	return &DocumentController{
		docService: docService,
	}
}

// CreateDocument handles document creation
func (dc *DocumentController) CreateDocument(c *gin.Context) {
	var doc models.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document format"})
		return
	}

	// Set creator ID from authenticated user
	userID := c.GetUint("userID")
	doc.CreatorID = userID

	if err := dc.docService.CreateDocument(c.Request.Context(), &doc); err != nil {
		// 根据错误类型返回不同的状态码和消息
		switch {
		case strings.Contains(err.Error(), "already exists"):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Document with this title already exists",
			})
		case strings.Contains(err.Error(), "service error"):
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create document",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, doc)
}

// UpdateDocument handles document updates
func (dc *DocumentController) UpdateDocument(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var doc models.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doc.ID = uint(id)
	doc.CreatorID = c.GetUint("userID")

	if err := dc.docService.UpdateDocument(c.Request.Context(), &doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, doc)
}

// DeleteDocument handles document deletion
func (dc *DocumentController) DeleteDocument(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	userID := c.GetUint("userID")
	if err := dc.docService.DeleteDocument(c.Request.Context(), uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetDocument retrieves a specific document
func (dc *DocumentController) GetDocument(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	doc, err := dc.docService.GetDocument(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	c.JSON(http.StatusOK, doc)
}

// ListDocuments retrieves a list of documents
func (dc *DocumentController) ListDocuments(c *gin.Context) {
	params := repository.DocumentListParams{
		Page:     1,
		PageSize: 10,
	}

	// Parse query parameters
	if page := c.Query("page"); page != "" {
		if pageNum, err := strconv.Atoi(page); err == nil {
			params.Page = pageNum
		}
	}

	if pageSize := c.Query("page_size"); pageSize != "" {
		if size, err := strconv.Atoi(pageSize); err == nil {
			params.PageSize = size
		}
	}

	if status := c.Query("status"); status != "" {
		params.Status = &status
	}

	if search := c.Query("search"); search != "" {
		params.Search = search
	}

	// Get creator ID from query if admin, otherwise use authenticated user ID
	if creatorID := c.Query("creator_id"); creatorID != "" {
		if id, err := strconv.ParseUint(creatorID, 10, 32); err == nil {
			uid := uint(id)
			params.CreatorID = &uid
		}
	}

	// Parse tags if provided
	if tags := c.QueryArray("tags"); len(tags) > 0 {
		params.Tags = tags
	}

	docs, total, err := dc.docService.ListDocuments(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"documents": docs,
		"total":     total,
		"page":      params.Page,
		"page_size": params.PageSize,
	})
}

// CreateVersion creates a new version of a document
func (dc *DocumentController) CreateVersion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	userID := c.GetUint("userID")
	if err := dc.docService.CreateVersion(c.Request.Context(), uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// GetVersions retrieves all versions of a document
func (dc *DocumentController) GetVersions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	versions, err := dc.docService.GetVersions(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, versions)
}

// ManageTags handles adding and removing tags from a document
func (dc *DocumentController) ManageTags(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var req struct {
		AddTags    []uint `json:"add_tags"`
		RemoveTags []uint `json:"remove_tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dc.docService.ManageTags(c.Request.Context(), uint(id), req.AddTags, req.RemoveTags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
