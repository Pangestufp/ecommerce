package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type SalesOrderDetailResponse struct {
	SalesOrder SalesOrderResponse           `json:"sales_order"`
	Details    []SalesOrderItemResponse     `json:"details"`
	Histories  []OrderStatusHistoryResponse `json:"histories"`
	Actions    []string                     `json:"actions"`
}

type SalesOrderResponse struct {
	SalesOrderID            string `json:"sales_order_id"`
	SalesOrderCode          string `json:"sales_order_code"`
	UserID                  string `json:"user_id"`
	CustomerName            string `json:"customer_name"`
	CustomerEmail           string `json:"customer_email"`
	CustomerPhone           string `json:"customer_phone"`
	CustomerAddressSnapshot string `json:"customer_address_snapshot"`
	OriginAddressSnapshot   string `json:"origin_address_snapshot"`

	TrackingNumber *string `json:"tracking_number"`

	ShippingCourier       string `json:"shipping_courier"`
	ShippingDisplayName   string `json:"shipping_display_name"`
	ShippingService       string `json:"shipping_service"`
	ShippingDescription   string `json:"shipping_description"`
	ShippingETD           string `json:"shipping_etd"`
	OriginDistrictID      string `json:"origin_district_id"`
	DestinationDistrictID string `json:"destination_district_id"`

	TotalWeight int `json:"total_weight"`

	ShippingFee       decimal.Decimal `json:"shipping_fee"`
	ShippingFeeFormat string          `json:"shipping_fee_format"`

	TotalBeforeDiscount       decimal.Decimal `json:"total_before_discount"`
	TotalBeforeDiscountFormat string          `json:"total_before_discount_format"`

	DiscountAmount       decimal.Decimal `json:"discount_amount"`
	DiscountAmountFormat string          `json:"discount_amount_format"`

	FinalTotal       decimal.Decimal `json:"final_total"`
	FinalTotalFormat string          `json:"final_total_format"`

	Status string  `json:"status"`
	Notes  *string `json:"notes"`

	ConfirmedAt *time.Time `json:"confirmed_at"`
	PaidAt      *time.Time `json:"paid_at"`
	DeliveredAt *time.Time `json:"delivered_at"`
	FinishAt    *time.Time `json:"finish_at"`
	CancelledAt *time.Time `json:"cancelled_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SalesOrderItemResponse struct {
	DetailID     string `json:"detail_id"`
	SalesOrderID string `json:"sales_order_id"`

	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	ProductCode string `json:"product_code"`
	ImageURL    string `json:"image_url"`

	Quantity    int `json:"quantity"`
	Weight      int `json:"weight"`
	TotalWeight int `json:"total_weight"`

	PriceBeforeDiscount       decimal.Decimal `json:"price_before_discount"`
	PriceBeforeDiscountFormat string          `json:"price_before_discount_format"`

	DiscountID *string `json:"discount_id"`

	DiscountAmount       decimal.Decimal `json:"discount_amount"`
	DiscountAmountFormat string          `json:"discount_amount_format"`

	FinalUnitPrice       decimal.Decimal `json:"final_unit_price"`
	FinalUnitPriceFormat string          `json:"final_unit_price_format"`

	FinalTotal       decimal.Decimal `json:"final_total"`
	FinalTotalFormat string          `json:"final_total_format"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderStatusHistoryResponse struct {
	HistoryID      string    `json:"history_id"`
	SalesOrderID   string    `json:"sales_order_id"`
	PreviousStatus *string   `json:"previous_status"`
	Status         string    `json:"status"`
	Note           string    `json:"note"`
	CreatedBy      string    `json:"created_by"`
	CreatedName    string    `json:"created_name"`
	CreatedAt      time.Time `json:"created_at"`
}
