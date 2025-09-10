package entity

import "gorm.io/gorm"

type Folder struct {
	gorm.Model
	Name         string `gorm:"not null"`
	TotalSheets  int    `gorm:"not null;default:480"`
	UsedSheets   int    `gorm:"not null;default:0"`
	FolderTypeID uint
	FolderType   FolderType
}
