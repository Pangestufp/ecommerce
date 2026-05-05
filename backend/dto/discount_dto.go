package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateDiscountRequest struct {
	ProductID     string  `json:"product_id"`
	DiscountName  string  `json:"discount_name"`
	DiscountType  string  `json:"discount_type"`
	DiscountValue float64 `json:"discount_value"`
	StartAt       string  `json:"start_at"`
	ExpiredAt     string  `json:"expired_at"`
}

type DiscountResponse struct {
	DiscountID           string          `json:"discount_id"`
	ProductID            string          `json:"product_id"`
	DiscountName         string          `json:"discount_name"`
	DiscountType         string          `json:"discount_type"`
	DiscountValue        decimal.Decimal `json:"discount_value"`
	DiscountValueFormat  string          `json:"discount_value_format"`
	DiscountAmountFormat string          `json:"discount_Amount_format"`
	StartAt              time.Time       `json:"start_at"`
	ExpiredAt            time.Time       `json:"expired_at"`
	FinalValue           string          `json:"final_value"`
	StartAtFormat        string          `json:"start_at_format"`
	ExpiredAtFormat      string          `json:"expired_at_format"`
	Status               int             `json:"status"`
	CreatedAt            time.Time       `json:"created_at"`
	StatusFormat         string          `json:"status_format"`
	CreatedBy            string          `json:"created_by"`
	CreatedName          string          `json:"created_name"`
}

type DiscountType struct {
	DiscountType string `json:"discount_type"`
}
