package dto

type ShippingCostRequest struct {
	Origin      string
	Destination string
	Weight      int
}

type ShippingOption struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Service     string `json:"service"`
	Description string `json:"description"`
	Cost        int    `json:"cost"`
	Etd         string `json:"etd"`
}
