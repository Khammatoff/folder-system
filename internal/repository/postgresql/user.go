package postgresql

import "folder-system/internal/entity"

func (r *Repository) CreateUser(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
