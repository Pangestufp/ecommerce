package repository

import (
	"backend/entity"
	"backend/errorhandler"
	"errors"

	"gorm.io/gorm"
)

type AddressRepository interface {
	CreateAddress(a *entity.UserAddress) error
	UpdateAddress(a *entity.UserAddress) error
	DeleteAddress(addressID string) error
	GetAddressByID(addressID string) (*entity.UserAddress, error)
	GetAddressByUserID(userID string) ([]entity.UserAddress, error)
	CountAddressByUserID(userID string) (int64, error)
	UnsetPrimaryByUserID(userID string) error
}

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) *addressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) CreateAddress(a *entity.UserAddress) error {
	return r.db.Create(a).Error
}

func (r *addressRepository) UpdateAddress(a *entity.UserAddress) error {
	result := r.db.Save(a)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return &errorhandler.InternalServerError{Message: "No Row Effect"}
	}
	return nil
}

func (r *addressRepository) DeleteAddress(addressID string) error {
	result := r.db.Delete(&entity.UserAddress{}, "address_id = ?", addressID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return &errorhandler.InternalServerError{Message: "No Row Effect"}
	}
	return nil
}

func (r *addressRepository) GetAddressByID(addressID string) (*entity.UserAddress, error) {
	var a entity.UserAddress
	err := r.db.First(&a, "address_id = ?", addressID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &errorhandler.NotFoundError{Message: "Address Not Found"}
		}
		return nil, err
	}
	return &a, nil
}

func (r *addressRepository) GetAddressByUserID(userID string) ([]entity.UserAddress, error) {
	var addresses []entity.UserAddress
	err := r.db.Where("user_id = ?", userID).Order("is_primary DESC, created_at DESC").Find(&addresses).Error
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *addressRepository) CountAddressByUserID(userID string) (int64, error) {
	var count int64
	err := r.db.Model(&entity.UserAddress{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *addressRepository) UnsetPrimaryByUserID(userID string) error {
	return r.db.Model(&entity.UserAddress{}).
		Where("user_id = ?", userID).
		Update("is_primary", 0).Error
}
