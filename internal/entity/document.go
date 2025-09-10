package entity

import "gorm.io/gorm"

type Document struct {
	gorm.Model
	Title          string `gorm:"not null"`
	SheetsCount    int    `gorm:"not null"`
	FolderID       *uint  // Use pointer to uint to allow NULL
	Folder         Folder
	DocumentTypeID uint
	DocumentType   DocumentType
}
