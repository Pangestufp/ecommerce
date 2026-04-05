package repository

import (
	"backend/entity"
	"backend/errorhandler"
	"errors"

	"gorm.io/gorm"
)

type TypeRepository interface {
	CreateType(t *entity.Type) error
	UpdateType(t *entity.Type) error
	DeleteType(typeID string) error
	GetTypeByID(typeID string) (*entity.Type, error)
	GetAllType() ([]entity.Type, error)
}

type typeRepository struct {
	db *gorm.DB
}

func NewTypeRepository(db *gorm.DB) *typeRepository {
	return &typeRepository{db: db}
}

func (r *typeRepository) CreateType(t *entity.Type) error {
	return r.db.Create(t).Error
}

func (r *typeRepository) UpdateType(t *entity.Type) error {
	result := r.db.Save(t)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return &errorhandler.InternalServerError{Message: "No Row Effect"}
	}
	return nil
}

func (r *typeRepository) DeleteType(typeID string) error {
	result := r.db.Model(&entity.Type{}).
		Where("type_id = ?", typeID).
		Updates(map[string]interface{}{"status": 0})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return &errorhandler.InternalServerError{Message: "No Row Effect"}
	}
	return nil
}

func (r *typeRepository) GetTypeByID(typeID string) (*entity.Type, error) {
	var t entity.Type
	err := r.db.First(&t, "type_id = ?", typeID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &errorhandler.NotFoundError{Message: "Type Not Found"}
		}
		return nil, err
	}
	return &t, nil
}

func (r *typeRepository) GetAllType() ([]entity.Type, error) {
	var types []entity.Type
	err := r.db.Where("status = ?", 1).Order("created_at ASC").Find(&types).Error
	if err != nil {
		return nil, err
	}
	return types, nil
}
