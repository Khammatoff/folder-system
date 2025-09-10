package postgresql

import (
	"fmt"
	"log"

	"folder-system/internal/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto Migrate - создает таблицы если их нет
	err = db.AutoMigrate(
		&entity.User{},
		&entity.Folder{},
		&entity.Document{},
		&entity.DocumentType{},
		&entity.FolderType{},
		&entity.FolderTypeAssignment{},
	)
	if err != nil {
		log.Printf("Warning: Auto migration completed with errors: %v", err)
	} else {
		log.Println("Database tables migrated successfully")
	}

	// Создаем начальные данные если таблицы пустые
	if err := createInitialData(db); err != nil {
		return nil, fmt.Errorf("failed to create initial data: %w", err)
	}

	return &Repository{db: db}, nil
}

func createInitialData(db *gorm.DB) error {
	// Создаем основные типы документов если их нет
	var docTypeCount int64
	if err := db.Model(&entity.DocumentType{}).Count(&docTypeCount).Error; err != nil {
		return err
	}

	if docTypeCount == 0 {
		documentTypes := []entity.DocumentType{
			{Name: "Contract"},     // ID = 1
			{Name: "Report"},       // ID = 2
			{Name: "Invoice"},      // ID = 3
			{Name: "Presentation"}, // ID = 4
		}
		if err := db.Create(&documentTypes).Error; err != nil {
			return err
		}
		log.Println("Created default document types")
	} else {
		log.Println("Document types already exist, skipping creation")
	}

	// Создаем основные типы папок если их нет
	var folderTypeCount int64
	if err := db.Model(&entity.FolderType{}).Count(&folderTypeCount).Error; err != nil {
		return err
	}

	if folderTypeCount == 0 {
		folderTypes := []entity.FolderType{
			{Name: "General"},
			{Name: "Financial"},
			{Name: "Legal"},
			{Name: "Technical"},
		}
		if err := db.Create(&folderTypes).Error; err != nil {
			return err
		}
		log.Println("Created default folder types")
	} else {
		log.Println("Folder types already exist, skipping creation")
	}

	// Создаем связи между типами если их нет
	var assignmentCount int64
	if err := db.Model(&entity.FolderTypeAssignment{}).Count(&assignmentCount).Error; err != nil {
		return err
	}

	if assignmentCount == 0 {
		assignments := []entity.FolderTypeAssignment{
			{DocumentTypeID: 1, FolderTypeID: 3}, // Contract -> Legal
			{DocumentTypeID: 2, FolderTypeID: 1}, // Report -> General
			{DocumentTypeID: 3, FolderTypeID: 2}, // Invoice -> Financial
			{DocumentTypeID: 4, FolderTypeID: 4}, // Presentation -> Technical
		}
		if err := db.Create(&assignments).Error; err != nil {
			return err
		}
		log.Println("Created default folder type assignments")
	} else {
		log.Println("Folder type assignments already exist, skipping creation")
	}

	return nil
}

func (r *Repository) DB() *gorm.DB {
	return r.db
}
