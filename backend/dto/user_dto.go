package dto

import "time"

type VUser struct {
	UserID     string    `json:"user_id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Status     int       `json:"status"`
	Role       string    `json:"role"`
	Address    string    `json:"address"`
	Phone      string    `json:"phone"`
	PostalCode string    `json:"postal_code"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
