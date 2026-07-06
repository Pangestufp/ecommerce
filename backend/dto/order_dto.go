package dto

type CreateOrderRequest struct {
	CheckoutID     string             `json:"checkout_id"`
	AddressID      string             `json:"address_id"`
	CourierCode    string             `json:"courier_code"`
	CourierService string             `json:"courier_service"`
	CourierName    string             `json:"courier_name"`
	ShippingCost   int                `json:"shipping_cost"`
	Note           *string            `json:"note"`
	Items          []OrderItemRequest `json:"items"`
	Subtotal       int64              `json:"subtotal"`
	Total          int64              `json:"total"`
}

type OrderItemRequest struct {
	ProductID  string  `json:"product_id"`
	Qty        int     `json:"qty"`
	DiscountID *string `json:"discount_id"`
	UnitPrice  int64   `json:"unit_price"`
}
