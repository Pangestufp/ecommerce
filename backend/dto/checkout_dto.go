package dto

import (
	"github.com/shopspring/decimal"
)

type VerifyCheckoutResponse struct {
	ProductPrice []ProductCheckoutData `json:"product_price"`
	User_Address []AddressResponse     `json:"user_address"`
}

type ProductCheckoutData struct {
	ProductID          string             `json:"product_id"`
	ProductPrice       decimal.Decimal    `json:"product_price"`
	ProductName        string             `json:"product_name"`
	ProductPriceFormat string             `json:"product_price_format"`
	Image              string             `json:"image"`
	AvailableStock     int                `json:"available_stock"`
	Qty                int                `json:"qty"`
	Discounts          []DiscountResponse `json:"discount"`
}
