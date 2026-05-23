package dto

import "time"

type CreateCourierRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type UpdateCourierRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type CourierResponse struct {
	ID              string    `json:"id"`
	Code            string    `json:"code"`
	Name            string    `json:"name"`
	Status          int       `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedAtFormat string    `json:"created_at_format"`
	UpdatedAtFormat string    `json:"updated_at_format"`
}
