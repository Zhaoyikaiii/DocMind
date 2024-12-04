package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Zhaoyikaiii/docmind/internal/models"
	"github.com/Zhaoyikaiii/docmind/internal/service"
	"github.com/Zhaoyikaiii/docmind/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UploadHandler struct {
	uploadDir   string
	maxSize     int64 // 最大文件大小（字节）
	fileService service.FileService
}

type UploadConfig struct {
	UploadDir string
	MaxSize   int64
}

func NewUploadHandler(config UploadConfig, fileService service.FileService) *UploadHandler {
	return &UploadHandler{
		uploadDir:   config.UploadDir,
		maxSize:     config.MaxSize,
		fileService: fileService,
	}
}

func (h *UploadHandler) HandleFileUpload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// 验证文件大小
	if header.Size > h.maxSize {
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("File too large. Maximum size is %d MB", h.maxSize/(1024*1024)),
		})
		return
	}

	// 验证文件类型
	if !h.isAllowedFileType(header.Filename) {
		c.JSON(400, gin.H{
			"error": "File type not allowed. Supported types: PDF, DOC, DOCX, TXT, MD",
		})
		return
	}

	// 创建上传目录
	uploadPath := filepath.Join(h.uploadDir, time.Now().Format("2006/01/02"))
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		utils.Logger.Error("Failed to create upload directory",
			zap.Error(err),
			zap.String("path", uploadPath))
		c.JSON(500, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// 生成存储文件名（使用UUID）
	storageName := h.generateStorageName(filepath.Ext(header.Filename))
	storagePath := filepath.Join(uploadPath, storageName)

	// 保存文件
	if err := h.saveFile(file, storagePath); err != nil {
		utils.Logger.Error("Failed to save file",
			zap.Error(err),
			zap.String("filename", header.Filename))
		c.JSON(500, gin.H{"error": "Failed to save file"})
		return
	}

	// 创建文件记录
	fileRecord := &models.File{
		OriginalName: header.Filename,
		StorageName:  storageName,
		Path:         storagePath,
		Size:         header.Size,
		ContentType:  header.Header.Get("Content-Type"),
		UploaderID:   c.GetUint("userID"),
	}

	if err := h.fileService.CreateFile(c.Request.Context(), fileRecord, file); err != nil {
		// 如果数据库保存失败，删除已上传的文件
		os.Remove(storagePath)
		c.JSON(500, gin.H{"error": "Failed to save file metadata"})
		return
	}

	c.JSON(200, gin.H{
		"id":            fileRecord.ID,
		"original_name": fileRecord.OriginalName,
		"size":          fileRecord.Size,
		"content_type":  fileRecord.ContentType,
		"created_at":    fileRecord.CreatedAt,
	})
}

func (h *UploadHandler) isAllowedFileType(filename string) bool {
	allowedTypes := map[string]bool{
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".txt":  true,
		".md":   true,
	}
	ext := strings.ToLower(filepath.Ext(filename))
	return allowedTypes[ext]
}

func (h *UploadHandler) generateStorageName(ext string) string {
	return fmt.Sprintf("%s%s", uuid.New().String(), ext)
}

func (h *UploadHandler) saveFile(file multipart.File, filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	return err
}
