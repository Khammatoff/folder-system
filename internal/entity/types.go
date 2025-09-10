package entity

import "gorm.io/gorm"

type DocumentType struct {
	gorm.Model
	Name string `gorm:"not null"`
}

type FolderType struct {
	gorm.Model
	Name string `gorm:"not null"`
}

// FolderTypeAssignment defines which folder types are acceptable for which document types.
type FolderTypeAssignment struct {
	gorm.Model
	DocumentTypeID uint `gorm:"not null"`
	FolderTypeID   uint `gorm:"not null"`
}
