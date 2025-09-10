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

func NewRepository(cfg interface{}) (*Repository, error) {
	// This function would typically take a config struct and build a DSN
	// For simplicity, we assume cfg is already a DSN string or GORM config
	var db *gorm.DB
	var err error

	if dsn, ok := cfg.(string); ok {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	} else {
		return nil, fmt.Errorf("invalid config for PostgreSQL repository")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto Migrate
	err = db.AutoMigrate(
		&entity.User{},
		&entity.Folder{},
		&entity.Document{},
		&entity.DocumentType{},
		&entity.FolderType{},
		&entity.FolderTypeAssignment{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}

	return &Repository{db: db}, nil
}

func (r *Repository) DB() *gorm.DB {
	return r.db
}
