package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateInventoryRequest struct {
	ProductID string  `json:"product_id" binding:"required"`
	CostPrice float64 `json:"cost_price" binding:"required"`
	Stock     int     `json:"stock" binding:"required"`
}

type UpdateInventoryRequest struct {
	CostPrice float64 `json:"cost_price"`
	Stock     int     `json:"stock"`
}

type InventoryResponse struct {
	BatchID         string          `json:"batch_id"`
	BatchCode       string          `json:"batch_code"`
	ProductID       string          `json:"product_id"`
	CostPrice       decimal.Decimal `json:"cost_price"`
	CostPriceFormat string          `json:"cost_price_format"`
	Stock           int             `json:"stock"`
	ReservedStock   int             `json:"reserved_stock"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}
