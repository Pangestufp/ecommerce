package dto

import (
	"time"

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

type CheckoutRedisData struct {
	CheckoutID string              `json:"checkout_id"`
	UserID     string              `json:"user_id"`
	Items      []CheckoutRedisItem `json:"items"`
	CreatedAt  time.Time           `json:"created_at"`
}

type CheckoutRedisItem struct {
	ProductID string `json:"product_id"`
	Qty       int    `json:"qty"`
}

type CreateCheckoutResponse struct {
	CheckoutID string `json:"checkout_id"`
}

type ShippingRequest struct {
	CheckoutID string `json:"checkout_id"`
	AddressID  string `json:"address_id"`
}

type CheckoutItemValidation struct {
	Product     ProductCheckoutData
	OriginalQty int
	AdjustedQty int
	WasAdjusted bool
	Removed     bool
}

type CheckoutValidationResult struct {
	RedisData CheckoutRedisData
	Items     []CheckoutItemValidation
}

type ConfirmCheckoutItemRequest struct {
	ProductID  string  `json:"product_id"`
	Qty        int     `json:"qty"`
	DiscountID *string `json:"discount_id"`
}

type ConfirmCheckoutRequest struct {
	CheckoutID          string                       `json:"checkout_id"`
	AddressID           string                       `json:"address_id"`
	CourierCode         string                       `json:"courier_code"`
	CourierService      string                       `json:"courier_service"`
	CourierName         string                       `json:"courier_name"`
	ShippingCost        decimal.Decimal              `json:"shipping_cost"`
	Subtotal            decimal.Decimal              `json:"subtotal"`
	Note                string                       `json:"note"`
	Items               []ConfirmCheckoutItemRequest `json:"items"`
	ShippingETD         string                       `json:"-"`
	ShippingDescription string                       `json:"-"`
}

type ConfirmCheckoutResponse struct {
	IdempotencyKey string `json:"idempotency_key"`
	Status         string `json:"status"` // selalu "QUEUED" saat pertama kali
}

type OrderItemResponse struct {
	ProductID      string          `json:"product_id"`
	ProductName    string          `json:"product_name"`
	ImageURL       string          `json:"image_url"`
	Quantity       int             `json:"quantity"`
	UnitPrice      decimal.Decimal `json:"unit_price"`
	DiscountAmount decimal.Decimal `json:"discount_amount"`
	FinalUnitPrice decimal.Decimal `json:"final_unit_price"`
	FinalTotal     decimal.Decimal `json:"final_total"`
}

type OrderSummaryResponse struct {
	SalesOrderID        string              `json:"sales_order_id"`
	SalesOrderCode      string              `json:"sales_order_code"`
	Status              string              `json:"status"`
	CustomerName        string              `json:"customer_name"`
	CustomerPhone       string              `json:"customer_phone"`
	ShippingAddress     string              `json:"shipping_address"`
	ShippingCourier     string              `json:"shipping_courier"`
	ShippingDisplayName string              `json:"shipping_display_name"`
	ShippingService     string              `json:"shipping_service"`
	ShippingFee         decimal.Decimal     `json:"shipping_fee"`
	TotalBeforeDiscount decimal.Decimal     `json:"total_before_discount"`
	DiscountAmount      decimal.Decimal     `json:"discount_amount"`
	FinalTotal          decimal.Decimal     `json:"final_total"`
	Notes               string              `json:"notes"`
	Items               []OrderItemResponse `json:"items"`
	ConfirmedAt         *time.Time          `json:"confirmed_at"`
	CreatedAt           time.Time           `json:"created_at"`
}

type CheckoutStatusResponse struct {
	IdempotencyKey string                `json:"idempotency_key"`
	Status         string                `json:"status"`
	ErrorMessage   string                `json:"error_message,omitempty"`
	Order          *OrderSummaryResponse `json:"order,omitempty"`
}
