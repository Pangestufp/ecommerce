package repository

import (
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"backend/helper"
	"log"

	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *entity.Product) error
	CreateProductImages(images []entity.ProductImage) error
	Update(product *entity.Product) error
	GetProductByID(productID string) (*entity.Product, error)
	GetProductImageByProductID(productID string) ([]entity.ProductImage, error)
	GetAll() ([]entity.Product, error)
	Delete(productID string) error
	DeleteImagesByProductID(productID string) ([]entity.ProductImage, error)
	GetProductByProductCode(productCode string) (*entity.Product, error)
	GetProductByProductSlug(productSlug string) (*entity.Product, error)
	GetProductEnriched(productID string) (*dto.ProductEnrichedForES, error)
	GetAllProductsPaginated(cursor *dto.Paginate, search string, limit int) ([]dto.ProductListRow, error)
	GetProductsEnrichedBatch(productIDs []string) ([]*dto.ProductEnrichedForES, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{db: db}
}

func (r *productRepository) CreateProduct(product *entity.Product) error {
	if err := r.db.Create(product).Error; err != nil {
		return err
	}

	return nil
}

func (r *productRepository) CreateProductImages(images []entity.ProductImage) error {

	for _, image := range images {
		if err := r.db.Create(&image).Error; err != nil {
			log.Printf("[ProductImage] Failed to save image %s: %v", image.ImageID, err)
		}
	}

	return nil
}

func (r *productRepository) Update(product *entity.Product) error {
	if err := r.db.Save(product).Error; err != nil {
		return err
	}
	return nil
}

func (r *productRepository) GetProductByID(productID string) (*entity.Product, error) {
	var product entity.Product

	if err := r.db.First(&product, "product_id = ?", productID).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) GetProductByProductCode(productCode string) (*entity.Product, error) {
	var product entity.Product

	if err := r.db.First(&product, "product_code = ?", productCode).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) GetProductByProductSlug(productSlug string) (*entity.Product, error) {
	var product entity.Product

	if err := r.db.First(&product, "product_slug = ?", productSlug).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) GetProductImageByProductID(productID string) ([]entity.ProductImage, error) {
	var images []entity.ProductImage
	r.db.Where("product_id = ?", productID).Find(&images)

	return images, nil
}

func (r *productRepository) GetAll() ([]entity.Product, error) {
	var products []entity.Product
	err := r.db.Where("status = ?", 1).Find(&products).Error
	return products, err
}

func (r *productRepository) Delete(productID string) error {

	result := r.db.Model(&entity.Product{}).
		Where("product_id = ?", productID).
		Updates(map[string]interface{}{
			"status": 0,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return &errorhandler.InternalServerError{Message: "No Row Effect"}
	}

	return nil
}

func (r *productRepository) DeleteImagesByProductID(productID string) ([]entity.ProductImage, error) {
	var images []entity.ProductImage

	if err := r.db.Where("product_id = ?", productID).Find(&images).Error; err != nil {
		return nil, err
	}

	r.db.Where("product_id = ?", productID).Delete(&entity.ProductImage{})

	return images, nil
}

func (r *productRepository) GetProductEnriched(productID string) (*dto.ProductEnrichedForES, error) {
	var product dto.ProductEnrichedForES

	now := helper.TimeNowWIB()

	query := `
		SELECT 
			p.product_id,
			p.product_code,
			p.product_name,
			p.product_slug,
			p.weight_gram,
			p.type_id,
			t.type_name,
			t.type_code,
			p.description,
			d.discount_id,
			d.discount_name,
			d.discount_type,
			d.discount_value,
			img.picture_path AS primary_image,
			img.image_id AS primary_image_id,
			COALESCE(pp.product_price, 0) AS product_price,
			CASE
				WHEN d.discount_type = 'amount' THEN d.discount_value
				WHEN d.discount_type = 'percentage' THEN COALESCE(pp.product_price, 0) * d.discount_value
				ELSE 0
			END AS best_discount,
			GREATEST(
				CASE
					WHEN d.discount_id IS NULL THEN pp.product_price
					WHEN d.discount_type = 'amount' THEN COALESCE(pp.product_price, 0) - d.discount_value
					WHEN d.discount_type = 'percentage' THEN COALESCE(pp.product_price, 0) - (COALESCE(pp.product_price, 0) * d.discount_value)
					ELSE 0
				END,
				1
			) AS best_price,
			SUM(COALESCE(i.stock, 0)) AS stock,
			SUM(COALESCE(i.reserved_stock, 0)) AS reserved_stock,
			SUM(COALESCE(i.stock, 0)) - SUM(COALESCE(i.reserved_stock, 0)) AS available_stock,
			CASE
				WHEN SUM(COALESCE(i.stock, 0)) = 0 OR COALESCE(pp.product_price, 0) = 0 THEN 0
				ELSE 1
			END AS available
		FROM products p
		JOIN types t ON p.type_id = t.type_id
		LEFT JOIN (
			SELECT
				product_id,
				SUM(COALESCE(stock, 0)) AS stock,
				SUM(COALESCE(reserved_stock, 0)) AS reserved_stock
			FROM inventories
			WHERE product_id = ?
			GROUP BY product_id
		) i ON p.product_id = i.product_id
		LEFT JOIN (
			SELECT
				product_id,
				discount_id,
				discount_name,
				discount_type,
				discount_value
			FROM discounts
			WHERE
				product_id = ?
				AND start_at <= ?
				AND expired_at >= ?
				AND status = 1
		) d ON p.product_id = d.product_id
		LEFT JOIN (
			SELECT product_price, product_id
			FROM product_prices
			WHERE product_id = ?
			ORDER BY created_at DESC
			LIMIT 1
		) pp ON p.product_id = pp.product_id
		LEFT JOIN (
			SELECT picture_path, product_id, image_id
			FROM product_images
			WHERE product_id = ?
			AND is_primary = 1
			LIMIT 1
		) img ON p.product_id = img.product_id
		WHERE p.product_id = ? AND p.status = 1
		GROUP BY 
			p.product_id, p.type_id, t.type_name, t.type_code,
			d.discount_id, d.discount_name, d.discount_type, d.discount_value,
			pp.product_price, img.image_id, img.picture_path
		ORDER BY best_discount DESC
		LIMIT 1
	`

	if err := r.db.Raw(query, productID, productID, now, now, productID, productID, productID).Scan(&product).Error; err != nil {
		return nil, err
	}

	product.BestDiscountFormat = helper.FormatRupiah(product.BestDiscount)
	product.BestPriceFormat = helper.FormatRupiah(product.BestPrice)
	product.ProductPriceFormat = helper.FormatRupiah(product.ProductPrice)

	return &product, nil
}

func (r *productRepository) GetAllProductsPaginated(cursor *dto.Paginate, search string, limit int) ([]dto.ProductListRow, error) {
	if cursor != nil {
		if cursor.Direction == nil {
			return nil, &errorhandler.BadRequestError{Message: "invalid cursor: direction is required"}
		}
		if *cursor.Direction == "prev" && (cursor.FirstID == nil || cursor.FirstCreatedAt == nil || *cursor.FirstID == "") {
			return nil, &errorhandler.BadRequestError{Message: "invalid cursor: FirstID and FirstCreatedAt required for prev direction"}
		}
	}

	now := helper.TimeNowWIB()

	if limit <= 0 {
		limit = 5
	}

	query := `
		SELECT
			p.product_id,
			p.product_code,
			p.product_name,
			p.product_slug,
			p.weight_gram,
			p.type_id,
			t.type_code || ' - ' || t.type_name AS type_name,
			p.description,
			p.status,
			p.created_at,
			p.updated_at,
			COALESCE(i.stock, 0) AS stock,
			COALESCE(i.reserved_stock, 0) AS reserved_stock,
			COALESCE(pp.product_price, 0) AS product_price,
			CASE
				WHEN pp.product_price IS NULL THEN 0
				ELSE 1
			END AS is_price_set,
			CASE
				WHEN i.stock IS NULL THEN 0
				ELSE 1
			END AS is_stock_set,
			COALESCE(d.available_discount, 0) AS available_discount
		FROM products p
		JOIN types t ON p.type_id = t.type_id
		LEFT JOIN (
			SELECT product_id,
				   SUM(stock) AS stock,
				   SUM(reserved_stock) AS reserved_stock
			FROM inventories
			GROUP BY product_id
		) i ON p.product_id = i.product_id
		LEFT JOIN (
			SELECT DISTINCT ON (product_id)
				product_id, product_price
			FROM product_prices
			ORDER BY product_id, created_at DESC
		) pp ON p.product_id = pp.product_id
		LEFT JOIN (
			SELECT product_id, COUNT(*) AS available_discount
			FROM discounts
			WHERE start_at <= ? AND expired_at >= ?
			AND status = 1
			GROUP BY product_id
		) d ON p.product_id = d.product_id
		 WHERE p.status = 1`

	args := []interface{}{now, now}

	if search != "" {
		search = "%" + search + "%"
	}

	if cursor != nil {
		if cursor.Direction != nil && *cursor.Direction == "prev" {
			if search != "" {
				query += ` AND (p.product_code ILIKE ? OR p.product_name ILIKE ?)`
				args = append(args, search, search)
			}
			query += ` AND (p.created_at, p.product_id) > (?, ?)`
			args = append(args, cursor.FirstCreatedAt, cursor.FirstID)
			query += ` ORDER BY p.created_at ASC, p.product_id ASC LIMIT ?`
		} else if cursor.Direction != nil && *cursor.Direction == "next" {
			if search != "" {
				query += ` AND (p.product_code ILIKE ? OR p.product_name ILIKE ?)`
				args = append(args, search, search)
			}
			query += ` AND (p.created_at, p.product_id) < (?, ?)`
			args = append(args, cursor.LastCreatedAt, cursor.LastID)
			query += ` ORDER BY p.created_at DESC, p.product_id DESC LIMIT ?`
		} else {
			if search != "" {
				query += ` AND (p.product_code ILIKE ? OR p.product_name ILIKE ?) ORDER BY p.created_at DESC, p.product_id DESC LIMIT ?`
				args = append(args, search, search)
			} else {
				query += ` ORDER BY p.created_at DESC, p.product_id DESC LIMIT ?`
			}
		}
	} else {
		if search != "" {
			query += ` AND (p.product_code ILIKE ? OR p.product_name ILIKE ?) ORDER BY p.created_at DESC, p.product_id DESC LIMIT ?`
			args = append(args, search, search)
		} else {
			query += ` ORDER BY p.created_at DESC, p.product_id DESC LIMIT ?`
		}
	}

	args = append(args, limit+1)

	var products []dto.ProductListRow
	if err := r.db.Raw(query, args...).Scan(&products).Error; err != nil {
		return nil, err
	}

	if cursor != nil && cursor.Direction != nil && *cursor.Direction == "prev" {
		for i, j := 0, len(products)-1; i < j; i, j = i+1, j-1 {
			products[i], products[j] = products[j], products[i]
		}
	}

	for i := range products {
		products[i].ProductPriceFormat = helper.FormatRupiah(products[i].ProductPrice)
	}

	return products, nil
}

func (r *productRepository) GetProductsEnrichedBatch(productIDs []string) ([]*dto.ProductEnrichedForES, error) {
	if len(productIDs) == 0 {
		return []*dto.ProductEnrichedForES{}, nil
	}

	var raw []*dto.ProductEnrichedForES
	now := helper.TimeNowWIB()

	query := `
		SELECT 
			p.product_id,
			p.product_code,
			p.product_name,
			p.product_slug,
			p.weight_gram,
			p.type_id,
			t.type_name,
			t.type_code,
			p.description,
			d.discount_id,
			d.discount_name,
			d.discount_type,
			d.discount_value,
			img.picture_path AS primary_image,
			img.image_id AS primary_image_id,
			COALESCE(pp.product_price, 0) AS product_price,
			CASE
				WHEN d.discount_type = 'amount' THEN d.discount_value
				WHEN d.discount_type = 'percentage' THEN COALESCE(pp.product_price, 0) * d.discount_value
				ELSE 0
			END AS best_discount,
			GREATEST(
				CASE
					WHEN d.discount_id IS NULL THEN COALESCE(pp.product_price, 0)
					WHEN d.discount_type = 'amount' THEN COALESCE(pp.product_price, 0) - d.discount_value
					WHEN d.discount_type = 'percentage' THEN COALESCE(pp.product_price, 0) - (COALESCE(pp.product_price, 0) * d.discount_value)
					ELSE 0
				END,
				1
			) AS best_price,
			SUM(COALESCE(i.stock, 0)) AS stock,
			SUM(COALESCE(i.reserved_stock, 0)) AS reserved_stock,
			SUM(COALESCE(i.stock, 0)) - SUM(COALESCE(i.reserved_stock, 0)) AS available_stock,
			CASE
				WHEN SUM(COALESCE(i.stock, 0)) = 0 OR COALESCE(pp.product_price, 0) = 0 THEN 0
				ELSE 1
			END AS available
		FROM products p
		JOIN types t ON p.type_id = t.type_id
		LEFT JOIN (
			SELECT
				product_id,
				SUM(COALESCE(stock, 0)) AS stock,
				SUM(COALESCE(reserved_stock, 0)) AS reserved_stock
			FROM inventories
			WHERE product_id IN (?)
			GROUP BY product_id
		) i ON p.product_id = i.product_id
		LEFT JOIN (
			SELECT
				product_id,
				discount_id,
				discount_name,
				discount_type,
				discount_value
			FROM discounts
			WHERE
				product_id IN (?)
				AND start_at <= ?
				AND expired_at >= ?
				AND status = 1
		) d ON p.product_id = d.product_id
		LEFT JOIN (
			SELECT product_price, product_id
			FROM (
				SELECT
					product_price,
					product_id,
					ROW_NUMBER() OVER (PARTITION BY product_id ORDER BY created_at DESC) AS rn
				FROM product_prices
				WHERE product_id IN (?)
			) ranked
			WHERE rn = 1
		) pp ON p.product_id = pp.product_id
		LEFT JOIN (
			SELECT picture_path, product_id, image_id
			FROM (
				SELECT
					picture_path,
					product_id,
					image_id,
					ROW_NUMBER() OVER (PARTITION BY product_id ORDER BY image_id) AS rn
				FROM product_images
				WHERE product_id IN (?)
				AND is_primary = 1
			) ranked
			WHERE rn = 1
		) img ON p.product_id = img.product_id
		WHERE p.product_id IN (?) AND p.status = 1
		GROUP BY
			p.product_id, p.type_id, t.type_name, t.type_code,
			d.discount_id, d.discount_name, d.discount_type, d.discount_value,
			pp.product_price, img.image_id, img.picture_path
		ORDER BY p.product_id, best_discount DESC
	`

	if err := r.db.Raw(query,
		productIDs, // inventories IN
		productIDs, // discounts IN
		now,        // start_at <=
		now,        // expired_at >=
		productIDs, // product_prices IN
		productIDs, // product_images IN
		productIDs, // WHERE p.product_id IN
	).Scan(&raw).Error; err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var products []*dto.ProductEnrichedForES

	for _, p := range raw {
		if !seen[p.ProductID] {
			seen[p.ProductID] = true
			p.BestDiscountFormat = helper.FormatRupiah(p.BestDiscount)
			p.BestPriceFormat = helper.FormatRupiah(p.BestPrice)
			p.ProductPriceFormat = helper.FormatRupiah(p.ProductPrice)
			products = append(products, p)
		}
	}

	return products, nil
}
