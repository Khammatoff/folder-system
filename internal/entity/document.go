package entity

import "gorm.io/gorm"

type Document struct {
	gorm.Model
	Title          string       `gorm:"not null" json:"title"`
	SheetsCount    int          `gorm:"not null" json:"sheets_count"`
	FolderID       *uint        `json:"folder_id"` // nullable
	Folder         *Folder      `gorm:"foreignKey:FolderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"folder,omitempty"`
	DocumentTypeID uint         `json:"document_type_id"`
	DocumentType   DocumentType `json:"document_type"`
}
