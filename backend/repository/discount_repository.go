package repository

import (
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"backend/helper"
	"errors"

	"gorm.io/gorm"
)

type DiscountRepository interface {
	Create(discount *entity.Discount) error
	Delete(discountID string) error
	GetByID(discountID string) (*entity.Discount, error)
	GetAllByProductID(productID string, cursor *dto.Paginate, search string, limit int) ([]entity.Discount, error)
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

func (r *discountRepository) GetAllByProductID(productID string, cursor *dto.Paginate, search string, limit int) ([]entity.Discount, error) {
	if limit <= 0 {
		limit = 5
	}

	var discounts []entity.Discount

	query := r.db.Model(&entity.Discount{}).Where("status = 1 AND product_id = ? AND expired_at >= ?", productID, helper.TimeNowWIB())

	if search != "" {
		query = query.Where("discount_name ILIKE ?", "%"+search+"%")
	}

	if cursor != nil {
		if cursor.Direction != nil && *cursor.Direction == "prev" {
			if cursor.FirstID != nil && cursor.FirstCreatedAt != nil {
				query = query.Where("(created_at, discount_id) > (?, ?)", cursor.FirstCreatedAt, cursor.FirstID)
				query = query.Order("created_at ASC, discount_id ASC")
			}
		} else {
			if cursor.Direction != nil && *cursor.Direction == "next" {
				query = query.Where("(created_at, discount_id) < (?, ?)", cursor.LastCreatedAt, cursor.LastID)
			}
			query = query.Order("created_at DESC, discount_id DESC")
		}
	} else {
		query = query.Order("created_at DESC, discount_id DESC")
	}

	err := query.Limit(limit + 1).Find(&discounts).Error
	if err != nil {
		return nil, err
	}

	if cursor != nil && cursor.Direction != nil && *cursor.Direction == "prev" {
		for i, j := 0, len(discounts)-1; i < j; i, j = i+1, j-1 {
			discounts[i], discounts[j] = discounts[j], discounts[i]
		}
	}

	return discounts, nil
}
