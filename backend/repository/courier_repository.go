package repository

import (
	"backend/entity"
	"backend/errorhandler"
	"errors"

	"gorm.io/gorm"
)

type CourierRepository interface {
	Create(courier *entity.Courier) error
	GetAll() ([]entity.Courier, error)
	GetActiveCouriers() ([]entity.Courier, error)
	GetByID(id string) (*entity.Courier, error)
	GetByCode(code string) (*entity.Courier, error)
	Update(courier *entity.Courier) error
	Toggle(id string) error
}

type courierRepository struct {
	db *gorm.DB
}

func NewCourierRepository(db *gorm.DB) *courierRepository {
	return &courierRepository{db: db}
}

func (r *courierRepository) Create(courier *entity.Courier) error {
	return r.db.Create(courier).Error
}

func (r *courierRepository) GetAll() ([]entity.Courier, error) {
	var couriers []entity.Courier
	err := r.db.Order("status DESC, name ASC").Find(&couriers).Error
	return couriers, err
}

func (r *courierRepository) GetActiveCouriers() ([]entity.Courier, error) {
	var couriers []entity.Courier
	err := r.db.Where("status = 1").Order("name ASC").Find(&couriers).Error
	return couriers, err
}

func (r *courierRepository) GetByID(id string) (*entity.Courier, error) {
	var courier entity.Courier
	err := r.db.First(&courier, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &errorhandler.NotFoundError{Message: "Kurir tidak ditemukan"}
		}
		return nil, err
	}
	return &courier, nil
}

func (r *courierRepository) GetByCode(code string) (*entity.Courier, error) {
	var courier entity.Courier
	err := r.db.First(&courier, "code = ?", code).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &errorhandler.NotFoundError{Message: "Kurir tidak ditemukan"}
		}
		return nil, err
	}
	return &courier, nil
}

func (r *courierRepository) Update(courier *entity.Courier) error {
	result := r.db.Save(courier)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return &errorhandler.InternalServerError{Message: "No Row Effect"}
	}
	return nil
}

func (r *courierRepository) Toggle(id string) error {
	// pakai CASE biar atomic, tidak perlu get dulu
	result := r.db.Model(&entity.Courier{}).
		Where("id = ?", id).
		Update("status", gorm.Expr("CASE WHEN status = 1 THEN 0 ELSE 1 END"))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return &errorhandler.NotFoundError{Message: "Kurir tidak ditemukan"}
	}
	return nil
}
