package models

import (
	"time"
	"gorm.io/gorm"
)

// File 文件模型
type File struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	OriginalName  string         `gorm:"size:255;not null" json:"original_name"`
	StorageName   string         `gorm:"size:255;not null;unique" json:"storage_name"`
	Path          string         `gorm:"size:512;not null" json:"path"`
	Size          int64          `gorm:"not null" json:"size"`
	ContentType   string         `gorm:"size:128" json:"content_type"`
	UploaderID    uint           `gorm:"not null" json:"uploader_id"`
	Uploader      User           `gorm:"foreignKey:UploaderID" json:"-"`
	DocumentID    *uint          `gorm:"index" json:"document_id,omitempty"`
	Document      *Document      `gorm:"foreignKey:DocumentID" json:"-"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
} 