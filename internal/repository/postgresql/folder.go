package postgresql

import "folder-system/internal/entity"

func (r *Repository) CreateFolder(folder *entity.Folder) error {
	return r.db.Create(folder).Error
}

func (r *Repository) GetFolderByID(id uint) (*entity.Folder, error) {
	var folder entity.Folder
	result := r.db.Preload("FolderType").First(&folder, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &folder, nil
}

func (r *Repository) UpdateFolder(folder *entity.Folder) error {
	return r.db.Save(folder).Error
}

func (r *Repository) FindFolderByTypeAndCapacity(folderTypeID uint, sheetsRequired int) (*entity.Folder, error) {
	var folder entity.Folder
	// Find the first folder of the given type that has enough free space.
	// (total_sheets - used_sheets) >= sheetsRequired
	result := r.db.
		Where("folder_type_id = ? AND (total_sheets - used_sheets) >= ?", folderTypeID, sheetsRequired).
		First(&folder)

	if result.Error != nil {
		return nil, result.Error
	}
	return &folder, nil
}
