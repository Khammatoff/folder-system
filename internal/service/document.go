package service

import (
	"errors"
	"folder-system/internal/entity"
	"folder-system/internal/repository"
)

type DocumentService interface {
	CreateDocument(title string, sheetsCount int, folderID *uint, docTypeID uint) (*entity.Document, error)
	GetDocument(id uint) (*entity.Document, error)
	UpdateDocument(id uint, title *string, sheetsCount *int, folderID *uint) (*entity.Document, error)
	DeleteDocument(id uint) error
}

type documentService struct {
	docRepo    repository.DocumentRepository
	folderRepo repository.FolderRepository
}

func NewDocumentService(docRepo repository.DocumentRepository, folderRepo repository.FolderRepository) DocumentService {
	return &documentService{docRepo: docRepo, folderRepo: folderRepo}
}

func (s *documentService) CreateDocument(title string, sheetsCount int, folderID *uint, docTypeID uint) (*entity.Document, error) {
	document := &entity.Document{
		Title:          title,
		SheetsCount:    sheetsCount,
		FolderID:       folderID,
		DocumentTypeID: docTypeID,
	}

	// If folder is specified, check capacity
	if folderID != nil {
		folder, err := s.folderRepo.GetFolderByID(*folderID)
		if err != nil {
			return nil, errors.New("folder not found")
		}

		if (folder.TotalSheets - folder.UsedSheets) < sheetsCount {
			return nil, errors.New("not enough space in the folder")
		}

		// Reserve space in the folder
		folder.UsedSheets += sheetsCount
		err = s.folderRepo.UpdateFolder(folder)
		if err != nil {
			return nil, errors.New("failed to update folder capacity")
		}
	}

	err := s.docRepo.CreateDocument(document)
	if err != nil {
		// If we reserved space but failed to create doc, we need to rollback?
		// This is a complex scenario. For production, use transactions.
		return nil, err
	}

	return document, nil
}

func (s *documentService) GetDocument(id uint) (*entity.Document, error) {
	return s.docRepo.GetDocumentByID(id)
}

func (s *documentService) UpdateDocument(id uint, title *string, sheetsCount *int, folderID *uint) (*entity.Document, error) {
	document, err := s.docRepo.GetDocumentByID(id)
	if err != nil {
		return nil, errors.New("document not found")
	}

	oldFolderID := document.FolderID
	oldSheetsCount := document.SheetsCount

	// Update simple fields
	if title != nil {
		document.Title = *title
	}
	if sheetsCount != nil {
		document.SheetsCount = *sheetsCount
	}

	// Handle folder change logic with transaction
	// For simplicity, we'll outline the logic without actual transaction handling here.
	// In production, you would use s.db.Transaction(...)

	// 1. If folder is being changed (to a new one or to nil)...
	if folderID != nil && (oldFolderID == nil || *folderID != *oldFolderID) {
		// ...free space in the old folder
		if oldFolderID != nil {
			oldFolder, err := s.folderRepo.GetFolderByID(*oldFolderID)
			if err == nil { // If old folder exists, decrement its used sheets
				oldFolder.UsedSheets -= oldSheetsCount
				s.folderRepo.UpdateFolder(oldFolder)
			}
		}

		// ...and reserve space in the new folder
		newFolder, err := s.folderRepo.GetFolderByID(*folderID)
		if err != nil {
			return nil, errors.New("new folder not found")
		}
		sheetsToUse := oldSheetsCount
		if sheetsCount != nil {
			sheetsToUse = *sheetsCount
		}
		if (newFolder.TotalSheets - newFolder.UsedSheets) < sheetsToUse {
			return nil, errors.New("not enough space in the new folder")
		}
		newFolder.UsedSheets += sheetsToUse
		err = s.folderRepo.UpdateFolder(newFolder)
		if err != nil {
			return nil, errors.New("failed to reserve space in new folder")
		}
		document.FolderID = folderID
	} else if folderID == nil && oldFolderID != nil {
		// If folder is being removed (set to nil), free space in the old folder
		oldFolder, err := s.folderRepo.GetFolderByID(*oldFolderID)
		if err == nil {
			oldFolder.UsedSheets -= oldSheetsCount
			s.folderRepo.UpdateFolder(oldFolder)
		}
		document.FolderID = nil
	} else if sheetsCount != nil && oldFolderID != nil {
		// If only sheet count changed and document is in a folder, adjust the reservation
		diff := *sheetsCount - oldSheetsCount
		folder, err := s.folderRepo.GetFolderByID(*oldFolderID)
		if err != nil {
			return nil, errors.New("folder not found")
		}
		if (folder.TotalSheets - folder.UsedSheets) < diff {
			return nil, errors.New("not enough space in the folder for the increase")
		}
		folder.UsedSheets += diff
		err = s.folderRepo.UpdateFolder(folder)
		if err != nil {
			return nil, errors.New("failed to adjust folder capacity")
		}
	}

	err = s.docRepo.UpdateDocument(document)
	if err != nil {
		return nil, err
	}

	return document, nil
}

func (s *documentService) DeleteDocument(id uint) error {
	// First, get the document to free up space in its folder
	document, err := s.docRepo.GetDocumentByID(id)
	if err != nil {
		return err
	}

	if document.FolderID != nil {
		folder, err := s.folderRepo.GetFolderByID(*document.FolderID)
		if err == nil {
			folder.UsedSheets -= document.SheetsCount
			s.folderRepo.UpdateFolder(folder)
		}
	}

	return s.docRepo.DeleteDocument(id)
}
