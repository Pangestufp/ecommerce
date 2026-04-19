package dto

import "time"

type PresignedURLRequest struct {
	Files []FileItem `json:"files"`
}

type FileItem struct {
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
}

type PresignedURLResponse struct {
	Uploads []UploadItem `json:"uploads"`
}

type UploadItem struct {
	UploadURL  string `json:"upload_url"`
	ObjectName string `json:"object_name"`
}

type CreateProductRequest struct {
	ProductCode string   `json:"product_code"`
	ProductName string   `json:"product_name"`
	WeightGram  int      `json:"weight_gram"`
	TypeID      string   `json:"type_id"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
}

type UpdateProductRequest struct {
	ProductName string   `json:"product_name"`
	ProductCode string   `json:"product_code"`
	WeightGram  int      `json:"weight_gram"`
	TypeID      string   `json:"type_id"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
}

type ProductResponse struct {
	ProductID   string                 `json:"product_id"`
	ProductCode string                 `json:"product_code"`
	ProductName string                 `json:"product_name"`
	ProductSlug string                 `json:"product_slug"`
	WeightGram  int                    `json:"weight_gram"`
	TypeID      string                 `json:"type_id"`
	Description string                 `json:"description"`
	Status      int                    `json:"status"`
	Images      []ProductImageResponse `json:"images"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type ProductImageResponse struct {
	ImageID     string `json:"image_id"`
	PicturePath string `json:"picture_path"`
	IsPrimary   int    `json:"is_primary"`
}

type ProductListRow struct {
	ProductID         string    `json:"product_id"`
	ProductCode       string    `json:"product_code"`
	ProductName       string    `json:"product_name"`
	ProductSlug       string    `json:"product_slug"`
	WeightGram        int       `json:"weight_gram"`
	TypeID            string    `json:"type_id"`
	TypeName          string    `json:"type_name"`
	Description       string    `json:"description"`
	Status            int       `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Stock             int64     `json:"stock"`
	ReservedStock     int64     `json:"reserved_stock"`
	ProductPrice      float64   `json:"product_price"`
	IsPriceSet        int       `json:"is_price_set"`
	IsStockSet        int       `json:"is_stock_set"`
	AvailableDiscount int64     `json:"available_discount"`
}
