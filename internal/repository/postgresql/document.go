package postgresql

import "folder-system/internal/entity"

func (r *Repository) CreateDocument(document *entity.Document) error {
	return r.db.Create(document).Error
}

func (r *Repository) GetDocumentByID(id uint) (*entity.Document, error) {
	var document entity.Document
	// Preload Folder and its type to check capacity later
	result := r.db.Preload("Folder").Preload("DocumentType").First(&document, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &document, nil
}

func (r *Repository) UpdateDocument(document *entity.Document) error {
	return r.db.Save(document).Error
}

func (r *Repository) DeleteDocument(id uint) error {
	return r.db.Delete(&entity.Document{}, id).Error
}
