package dto

type TypeRequest struct {
	TypeCode string `json:"type_code"`
	TypeName string `json:"type_name"`
	TypeDesc string `json:"type_desc"`
}

type TypeResponse struct {
	TypeID   string `json:"type_id"`
	TypeCode string `json:"type_code"`
	TypeName string `json:"type_name"`
	TypeDesc string `json:"type_desc"`
}
