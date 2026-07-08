package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type SalesOrder struct {
	SalesOrderID            string `gorm:"primaryKey"`
	SalesOrderCode          string
	UserID                  string
	CustomerName            string
	CustomerEmail           string
	CustomerPhone           string
	CustomerAddressSnapshot string
	OriginAddressSnapshot   string
	TrackingNumber          *string
	ShippingCourier         string
	ShippingDisplayName     string
	ShippingService         string
	ShippingDescription     string
	ShippingFee             decimal.Decimal
	ShippingETD             string
	OriginDistrictID        string
	DestinationDistrictID   string
	TotalWeight             int
	TotalBeforeDiscount     decimal.Decimal
	DiscountAmount          decimal.Decimal
	FinalTotal              decimal.Decimal
	Status                  string
	Notes                   *string
	ConfirmedAt             *time.Time
	PaidAt                  *time.Time
	DeliveredAt             *time.Time
	FinishAt                *time.Time
	CancelledAt             *time.Time
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type SalesOrderDetail struct {
	DetailID            string `gorm:"primaryKey"`
	SalesOrderID        string
	ProductID           string
	ImageURL            string
	ProductName         string
	ProductCode         string
	Quantity            int
	Weight              int
	TotalWeight         int
	PriceBeforeDiscount decimal.Decimal
	DiscountID          *string
	DiscountAmount      decimal.Decimal
	FinalUnitPrice      decimal.Decimal
	FinalTotal          decimal.Decimal
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
