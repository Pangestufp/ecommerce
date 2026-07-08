package entity

import "time"

type OrderStatusHistory struct {
	HistoryID      string `gorm:"primaryKey"`
	SalesOrderID   string
	PreviousStatus *string
	Status         string
	Note           string
	CreatedAt      time.Time
	CreatedBy      string
	CreatedName    string
}
