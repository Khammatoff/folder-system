package service

import (
	"folder-system/internal/entity"
	"folder-system/internal/repository"
)

type FolderService interface {
	GetRecommendedFolder(documentTypeID uint, sheetsCount int) (*entity.Folder, error)
}

type folderService struct {
	folderRepo repository.FolderRepository
}

func NewFolderService(folderRepo repository.FolderRepository) FolderService {
	return &folderService{folderRepo: folderRepo}
}

func (s *folderService) GetRecommendedFolder(documentTypeID uint, sheetsCount int) (*entity.Folder, error) {
	// This is a simplified recommendation.
	// In a real application, you would join through FolderTypeAssignment to find a folder
	// of a type that is allowed for this document type, and that has enough space.
	// For this example, we'll assume the documentTypeID can be directly used to find a folder type.
	// You would need to implement the logic to map document type to folder type.

	// Placeholder: Assume documentTypeID corresponds to a folder type ID (1:1 for simplicity)
	// In a real app, you'd have a separate method to get folder type from document type.
	folderTypeID := documentTypeID

	folder, err := s.folderRepo.FindFolderByTypeAndCapacity(folderTypeID, sheetsCount)
	if err != nil {
		return nil, err // "not found" is a valid outcome here, not necessarily an error
	}
	return folder, nil
}
