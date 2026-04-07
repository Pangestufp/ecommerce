package entity

import "time"

type Transaction struct {
	TransactionID string
	BatchID       string
	Type          string
	Quantity      int
	ReferenceType string
	ReferenceID   string
	Note          string
	CreatedAt     time.Time
}
