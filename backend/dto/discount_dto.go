package dto

import "time"

type CreateDiscountRequest struct {
	ProductID     string    `json:"product_id"`
	DiscountName  string    `json:"discount_name"`
	DiscountType  string    `json:"discount_type"`
	DiscountValue float64   `json:"discount_value"`
	StartAt       time.Time `json:"start_at"`
	ExpiredAt     time.Time `json:"expired_at"`
}

type DiscountResponse struct {
	DiscountID    string    `json:"discount_id"`
	ProductID     string    `json:"product_id"`
	DiscountName  string    `json:"discount_name"`
	DiscountType  string    `json:"discount_type"`
	DiscountValue float64   `json:"discount_value"`
	StartAt       time.Time `json:"start_at"`
	ExpiredAt     time.Time `json:"expired_at"`
	Status        int       `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}
