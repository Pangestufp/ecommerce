package repository

import (
	"backend/dto"
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
	GetAllTypePaginate(cursor *dto.Paginate, limit int) ([]entity.Type, error)
	GetAllType() ([]entity.Type, error)
	GetTypeByTypeCode(typeCode string) (*entity.Type, error)
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

func (r *typeRepository) GetTypeByTypeCode(typeCode string) (*entity.Type, error) {
	var t entity.Type
	err := r.db.First(&t, "type_code = ?", typeCode).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &errorhandler.NotFoundError{Message: "Type Not Found"}
		}
		return nil, err
	}
	return &t, nil
}

func (r *typeRepository) GetAllTypePaginate(cursor *dto.Paginate, limit int) ([]entity.Type, error) {
	if limit <= 0 {
		limit = 5
	}

	var types []entity.Type

	query := r.db.Model(&entity.Type{}).
		Where("status = ?", 1)

	if cursor != nil && cursor.LastID != nil && cursor.LastCreatedAt != nil {
		query = query.Where(
			"(created_at, type_id) < (?, ?)",
			cursor.LastCreatedAt,
			cursor.LastID,
		)
	}

	err := query.
		Order("created_at DESC, type_id DESC").
		Limit(limit).
		Find(&types).Error

	if err != nil {
		return nil, err
	}

	return types, nil
}

func (r *typeRepository) GetAllType() ([]entity.Type, error) {

	var types []entity.Type

	err := r.db.Where("status = ?", 1).Order("created_at DESC").Find(&types).Error
	if err != nil {
		return nil, err
	}

	return types, nil
}
