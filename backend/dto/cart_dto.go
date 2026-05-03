package dto

import "github.com/shopspring/decimal"

type CartRequest struct {
	ListCart []ListCart `json:"list_cart"`
}

type ListCart struct {
	ProductID string `json:"product_id"`
	Qty       int    `json:"qty"`
}

type CartResponse struct {
	ListProduct []CartVerifiedProduct `json:"list_product"`
	ListSave    []NewSave             `json:"list_save"`
	IsNote      int                   `json:"is_note"`
	Note        string                `json:"note"`
	TotalNow    decimal.Decimal       `json:"total_now"`
}

type CartVerifiedProduct struct {
	ProductID      string          `json:"product_id"`
	ProductName    string          `json:"product_name"`
	Image          string          `json:"image"`
	IsAvailable    int             `json:"is_available"`
	AvailableStock int             `json:"available_stock"`
	Qty            int             `json:"qty"`
	BestPrice      decimal.Decimal `json:"Price"`
	PriceFormat    string          `json:"Price_format"`
}

type NewSave struct {
	ID          string `json:"id"`
	ProductName string `json:"product_name"`
	Qty         int    `json:"qty"`
	Image       string `json:"image"`
}
