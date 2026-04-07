package repository

import (
	"backend/entity"
	"backend/errorhandler"
	"errors"

	"gorm.io/gorm"
)

type ProductPriceRepository interface {
	Create(price *entity.ProductPrice) error
	GetLatestByProductID(productID string) (*entity.ProductPrice, error)
	GetAllByProductID(productID string) ([]entity.ProductPrice, error)
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

func (r *productPriceRepository) GetAllByProductID(productID string) ([]entity.ProductPrice, error) {
	var prices []entity.ProductPrice
	err := r.db.Where("product_id = ?", productID).
		Order("created_at DESC").
		Find(&prices).Error
	if err != nil {
		return nil, err
	}
	return prices, nil
}
