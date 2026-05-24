package dto

import "time"

type TransactionResponse struct {
	TransactionID string    `json:"transaction_id"`
	BatchID       string    `json:"batch_id"`
	Type          string    `json:"type"` 
	Quantity      int       `json:"quantity"`
	ReferenceType string    `json:"reference_type"`
	ReferenceID   string    `json:"reference_id"`
	Note          string    `json:"note"`
	CreatedAt     time.Time `json:"created_at"`
}