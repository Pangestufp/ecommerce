package repository

import (
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"errors"

	"gorm.io/gorm"
)

type ProductPriceRepository interface {
	Create(price *entity.ProductPrice) error
	GetLatestByProductID(productID string) (*entity.ProductPrice, error)
	GetAllByProductID(productID string, cursor *dto.Paginate, limit int) ([]entity.ProductPrice, error)
}

type productPriceRepository struct {
	db *gorm.DB
}

func NewProductPriceRepository(db *gorm.DB) *productPriceRepository {
	return &productPriceRepository{db: db}
}

func (r *productPriceRepository) Create(price *entity.ProductPrice) error {
	return r.db.Create(price).Error
}

func (r *productPriceRepository) GetLatestByProductID(productID string) (*entity.ProductPrice, error) {
	var price entity.ProductPrice
	err := r.db.Where("product_id = ?", productID).
		Order("created_at DESC").
		First(&price).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &errorhandler.NotFoundError{Message: "Price Not Found"}
		}
		return nil, err
	}
	return &price, nil
}

func (r *productPriceRepository) GetAllByProductID(productID string, cursor *dto.Paginate, limit int) ([]entity.ProductPrice, error) {
	if limit <= 0 {
		limit = 5
	}

	var prices []entity.ProductPrice

	query := r.db.Model(&entity.ProductPrice{}).
		Where("product_id = ?", productID)

	if cursor != nil {
		if cursor.Direction != nil && *cursor.Direction == "prev" {
			if cursor.FirstID != nil && cursor.FirstCreatedAt != nil {
				query = query.Where("(created_at, price_id) > (?, ?)", cursor.FirstCreatedAt, cursor.FirstID)
				query = query.Order("created_at ASC, price_id ASC")
			}
		} else {
			if cursor.Direction != nil && *cursor.Direction == "next" {
				query = query.Where("(created_at, price_id) < (?, ?)", cursor.LastCreatedAt, cursor.LastID)
			}
			query = query.Order("created_at DESC, price_id DESC")
		}
	} else {
		query = query.Order("created_at DESC, price_id DESC")
	}

	err := query.Limit(limit + 1).Find(&prices).Error
	if err != nil {
		return nil, err
	}

	if cursor != nil && cursor.Direction != nil && *cursor.Direction == "prev" {
		for i, j := 0, len(prices)-1; i < j; i, j = i+1, j-1 {
			prices[i], prices[j] = prices[j], prices[i]
		}
	}

	return prices, nil
}
