package repository

import "folder-system/internal/entity"

// UserRepository defines the interface for user data access.
type UserRepository interface {
	CreateUser(user *entity.User) error
	GetUserByEmail(email string) (*entity.User, error)
}

// FolderRepository defines the interface for folder data access.
type FolderRepository interface {
	CreateFolder(folder *entity.Folder) error
	GetFolderByID(id uint) (*entity.Folder, error)
	UpdateFolder(folder *entity.Folder) error
	FindFolderByTypeAndCapacity(folderTypeID uint, sheetsRequired int) (*entity.Folder, error)
}

// DocumentRepository defines the interface for document data access.
type DocumentRepository interface {
	CreateDocument(document *entity.Document) error
	GetDocumentByID(id uint) (*entity.Document, error)
	UpdateDocument(document *entity.Document) error
	DeleteDocument(id uint) error
}
