package dto

import "time"

type ProductEnrichedForES struct {
	ProductID          string
	ProductCode        string
	ProductName        string
	ProductSlug        string
	WeightGram         int
	TypeID             string
	TypeName           string
	TypeCode           string
	Description        string
	DiscountID         *string
	DiscountName       *string
	DiscountType       *string
	DiscountValue      *float64
	ProductPrice       float64
	BestDiscount       float64
	BestPrice          float64
	ProductPriceFormat string `gorm:"-"`
	BestDiscountFormat string `gorm:"-"`
	BestPriceFormat    string `gorm:"-"`
	Stock              int64
	ReservedStock      int64
	AvailableStock     int64
	Available          int
	Images             []ProductImageResponse `gorm:"-"`
	Discounts          []DiscountResponse     `gorm:"-"`
}

type ProductEvent struct {
	ProductID string
	Type      string
	Timestamp time.Time
}
