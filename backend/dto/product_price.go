package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateProductPriceRequest struct {
	ProductID    string  `json:"product_id"`
	ProductPrice float64 `json:"product_price"`
}

type ProductPriceResponse struct {
	PriceID            string          `json:"price_id"`
	ProductID          string          `json:"product_id"`
	ProductPrice       decimal.Decimal `json:"product_price"`
	ProductPriceFormat string          `json:"product_price_format"`
	CreatedAt          time.Time       `json:"created_at"`
	CreatedBy          string          `json:"created_by"`
	CreatedName        string          `json:"created_name"`
}
