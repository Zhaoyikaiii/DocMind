package controllers

import (
	"net/http"
	"strconv"

	"github.com/Zhaoyikaiii/docmind/internal/models"
	"github.com/Zhaoyikaiii/docmind/internal/service"
	"github.com/gin-gonic/gin"
)

type FileController struct {
	fileService service.FileService
}

func NewFileController(fileService service.FileService) *FileController {
	return &FileController{
		fileService: fileService,
	}
}

// GetFile 获取文件信息
func (fc *FileController) GetFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	file, err := fc.fileService.GetFile(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.JSON(http.StatusOK, file)
}

// ListFiles 获取文件列表
func (fc *FileController) ListFiles(c *gin.Context) {
	var params models.FileListParams
	params.Page = 1
	params.PageSize = 10

	// 解析查询参数
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

	// 处理其他过滤条件
	if contentType := c.Query("content_type"); contentType != "" {
		params.ContentType = &contentType
	}

	if search := c.Query("search"); search != "" {
		params.Search = search
	}

	// 获取上传者ID
	if uploaderID := c.Query("uploader_id"); uploaderID != "" {
		if id, err := strconv.ParseUint(uploaderID, 10, 32); err == nil {
			uid := uint(id)
			params.UploaderID = &uid
		}
	}

	// 获取关联文档ID
	if documentID := c.Query("document_id"); documentID != "" {
		if id, err := strconv.ParseUint(documentID, 10, 32); err == nil {
			did := uint(id)
			params.DocumentID = &did
		}
	}

	files, total, err := fc.fileService.ListFiles(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files":     files,
		"total":     total,
		"page":      params.Page,
		"page_size": params.PageSize,
	})
}

// DeleteFile 删除文件
func (fc *FileController) DeleteFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	userID := c.GetUint("userID")
	if err := fc.fileService.DeleteFile(c.Request.Context(), uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	c.Status(http.StatusNoContent)
}

// AssociateWithDocument 将文件关联到文档
func (fc *FileController) AssociateWithDocument(c *gin.Context) {
	fileID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	var req struct {
		DocumentID uint `json:"document_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	file := &models.File{
		ID:         uint(fileID),
		DocumentID: &req.DocumentID,
		UploaderID: c.GetUint("userID"),
	}

	if err := fc.fileService.UpdateFile(c.Request.Context(), file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate file with document"})
		return
	}

	c.Status(http.StatusOK)
} 