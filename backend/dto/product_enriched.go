package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type ProductEnrichedForES struct {
	ProductID          string                 `json:"product_id"`
	ProductCode        string                 `json:"product_code"`
	ProductName        string                 `json:"product_name"`
	ProductSlug        string                 `json:"product_slug"`
	WeightGram         int                    `json:"weight_gram"`
	TypeID             string                 `json:"type_id"`
	TypeName           string                 `json:"type_name"`
	TypeCode           string                 `json:"type_code"`
	Description        string                 `json:"description"`
	DiscountID         *string                `json:"discount_id"`
	DiscountName       *string                `json:"discount_name"`
	DiscountType       *string                `json:"discount_type"`
	DiscountValue      *decimal.Decimal       `json:"discount_value"`
	ProductPrice       decimal.Decimal        `json:"product_price"`
	BestDiscount       decimal.Decimal        `json:"best_discount"`
	BestPrice          decimal.Decimal        `json:"best_price"`
	ProductPriceFormat string                 `json:"product_price_format" gorm:"-"`
	BestDiscountFormat string                 `json:"best_discount_format" gorm:"-"`
	BestPriceFormat    string                 `json:"best_price_format" gorm:"-"`
	Stock              int64                  `json:"stock"`
	ReservedStock      int64                  `json:"reserved_stock"`
	AvailableStock     int64                  `json:"available_stock"`
	Available          int                    `json:"available"`
	PrimaryImage       string                 `json:"primary_image"`
	PrimaryImageID     string                 `json:"primary_image_id"`
	Images             []ProductImageResponse `json:"images" gorm:"-"`
	Discounts          []DiscountResponse     `json:"discounts" gorm:"-"`
}

type ProductEvent struct {
	ProductID string    `json:"product_id"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}
