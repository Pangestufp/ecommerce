package entity

import "time"

type StockReservation struct {
	ReservedID   string `gorm:"primaryKey"`
	BatchID      string
	SalesOrderID string
	DetailID     string
	Quantity     int
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
