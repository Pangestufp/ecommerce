package dto

import "time"

type LogResponse struct {
	LogID         string    `json:"log_id"`
	ReferenceType string    `json:"reference_type"`
	ReferenceID   string    `json:"reference_id"`
	ReferenceName string    `json:"reference_name"`
	Note          string    `json:"note"`
	SourceID      string    `json:"source_id"`   
	SourceName    string    `json:"source_name"` 
	SourceType    string    `json:"source_type"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
	CreatedName   string    `json:"created_name"`
}