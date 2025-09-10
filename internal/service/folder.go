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
	folderTypeID := documentTypeID

	folder, err := s.folderRepo.FindFolderByTypeAndCapacity(folderTypeID, sheetsCount)
	if err != nil {
		return nil, err // "not found" is a valid outcome here, not necessarily an error
	}
	return folder, nil
}
