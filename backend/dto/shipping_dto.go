package dto

type ShippingCostRequest struct {
	Origin          string
	Destination     string
	OriginName      string
	DestinationName string
	Weight          int
}

type ShippingResponse struct {
	OriginName      string                 `json:"origin_name"`
	DestinationName string                 `json:"destination_name"`
	ShippingService []ShippingGroupService `json:"shipping_service"`
}

type ShippingGroupService struct {
	Name   string           `json:"name"`
	Code   string           `json:"code"`
	Option []ShippingOption `json:"option"`
}

type ShippingOption struct {
	Service       string `json:"service"`
	Description   string `json:"description"`
	Cost          int    `json:"cost"`
	Etd           string `json:"etd"`
	DisplayName   string `json:"display_name"`
	Group         string `json:"group"`
	IsRecommended bool   `json:"is_recommended"`
}
