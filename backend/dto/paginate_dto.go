package dto

type Paginate struct {
	Total  int `json:"total"`
	LastID int `json:"last_id"`
}
