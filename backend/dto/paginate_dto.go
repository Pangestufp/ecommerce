package dto

import "time"

type Paginate struct {
	FirstID        *string    `json:"first_id"`
	FirstCreatedAt *time.Time `json:"first_created_at"`
	LastID         *string    `json:"last_id"`
	LastCreatedAt  *time.Time `json:"last_created_at"`
	HasNext        *string    `json:"has_next"`
	HasPrev        *string    `json:"has_prev"`
	Direction      *string    `json:"direction"`
}
