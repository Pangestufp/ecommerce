package repository

import (
	"backend/entity"
	"backend/errorhandler"
	"errors"

	"gorm.io/gorm"
)

type DiscountRepository interface {
	Create(discount *entity.Discount) error
	Delete(discountID string) error
	GetByID(discountID string) (*entity.Discount, error)
}

type discountRepository struct {
	db *gorm.DB
}

func NewDiscountRepository(db *gorm.DB) *discountRepository {
	return &discountRepository{db: db}
}

func (r *discountRepository) Create(discount *entity.Discount) error {
	return r.db.Create(discount).Error
}

func (r *discountRepository) Delete(discountID string) error {
	result := r.db.Model(&entity.Discount{}).
		Where("discount_id = ?", discountID).
		Update("status", 0)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return &errorhandler.InternalServerError{Message: "No Row Effect"}
	}
	return nil
}

func (r *discountRepository) GetByID(discountID string) (*entity.Discount, error) {
	var discount entity.Discount
	err := r.db.First(&discount, "discount_id = ?", discountID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &errorhandler.NotFoundError{Message: "Discount Not Found"}
		}
		return nil, err
	}
	return &discount, nil
}
