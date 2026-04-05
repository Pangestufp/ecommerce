package repository

import (
	"backend/entity"
	"backend/errorhandler"
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
