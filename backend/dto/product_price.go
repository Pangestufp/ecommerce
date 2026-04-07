package dto

import "time"

type CreateProductPriceRequest struct {
	ProductID    string  `json:"product_id" binding:"required"`
	ProductPrice float64 `json:"product_price" binding:"required"`
}

type ProductPriceResponse struct {
	PriceID      string    `json:"price_id"`
	ProductID    string    `json:"product_id"`
	ProductPrice float64   `json:"product_price"`
	CreatedAt    time.Time `json:"created_at"`
}
