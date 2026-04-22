package repository

import (
	"backend/dto"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(userID string) (*dto.VUser, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUserByID(userID string) (*dto.VUser, error) {
	var user dto.VUser

	err := r.db.First(&user, "user_id = ?", userID).Error

	return &user, err
}
