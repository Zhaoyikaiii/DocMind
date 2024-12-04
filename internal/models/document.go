package models

import (
	"time"

	"gorm.io/gorm"
)

type Document struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Title     string         `gorm:"size:255;not null;uniqueIndex:idx_title_creator" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	Version   int            `gorm:"default:1" json:"version"`
	Status    string         `gorm:"size:20;default:'draft'" json:"status"` // draft, published, archived
	CreatorID uint           `gorm:"not null;uniqueIndex:idx_title_creator" json:"creator_id"`
	Creator   User           `gorm:"foreignKey:CreatorID" json:"creator"`
	ParentID  *uint          `gorm:"default:null" json:"parent_id"`
	Path      string         `gorm:"size:255" json:"path"`
	Tags      []Tag          `gorm:"many2many:document_tags;" json:"tags"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type Tag struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:50;not null;unique" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// DocumentVersion records the version history of a document
type DocumentVersion struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	DocumentID uint      `gorm:"not null" json:"document_id"`
	Version    int       `gorm:"not null" json:"version"`
	Title      string    `gorm:"size:255;not null" json:"title"`
	Content    string    `gorm:"type:text" json:"content"`
	CreatedBy  uint      `gorm:"not null" json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
}
