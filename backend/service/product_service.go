package service

import (
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"backend/helper"
	"backend/repository"
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

type ProductService interface {
	GeneratePresignedURLs(req dto.PresignedURLRequest) (*dto.PresignedURLResponse, error)
	Create(req dto.CreateProductRequest) (*dto.ProductResponse, error)
	Update(productID string, req dto.UpdateProductRequest) (*dto.ProductResponse, error)
	GetByID(productID string) (*dto.ProductResponse, error)
	GetAll() ([]dto.ProductResponse, error)
	Delete(productID string) error
}

type productService struct {
	repo   repository.ProductRepository
	minio  *minio.Client
	redis  *redis.Client
	bucket string
}

func NewProductService(repo repository.ProductRepository, minio *minio.Client, redis *redis.Client, bucket string) *productService {
	return &productService{repo: repo, minio: minio, redis: redis, bucket: bucket}
}

func (s *productService) GeneratePresignedURLs(req dto.PresignedURLRequest) (*dto.PresignedURLResponse, error) {
	var uploads []dto.UploadItem

	for _, file := range req.Files {
		ext := filepath.Ext(file.FileName)
		objectName := "temp/" + uuid.New().String() + ext

		url, err := s.minio.PresignedPutObject(
			context.Background(),
			s.bucket,
			objectName,
			15*time.Minute,
		)
		if err != nil {
			return nil, &errorhandler.InternalServerError{Message: err.Error()}
		}

		uploads = append(uploads, dto.UploadItem{
			UploadURL:  url.String(),
			ObjectName: objectName,
		})
	}

	return &dto.PresignedURLResponse{Uploads: uploads}, nil
}

func (s *productService) moveFromTemp(objectName string) (string, error) {
	permanentName := "products/" + uuid.New().String() + filepath.Ext(objectName)

	src := minio.CopySrcOptions{Bucket: s.bucket, Object: objectName}
	dst := minio.CopyDestOptions{Bucket: s.bucket, Object: permanentName}

	if _, err := s.minio.CopyObject(context.Background(), dst, src); err != nil {
		return "", err
	}

	s.minio.RemoveObject(context.Background(), s.bucket, objectName, minio.RemoveObjectOptions{})

	return permanentName, nil
}

func (s *productService) buildImages(productID string, objectNames []string) ([]entity.ProductImage, []error) {
	var images []entity.ProductImage
	var errs []error
	primarySet := false

	for _, objectName := range objectNames {
		permanentName, err := s.moveFromTemp(objectName)
		if err != nil {
			errs = append(errs, fmt.Errorf("object %s: %w", objectName, err))
			continue
		}

		primary := 0
		if !primarySet {
			primary = 1
			primarySet = true
		}

		images = append(images, entity.ProductImage{
			ImageID:     uuid.New().String(),
			ProductID:   productID,
			PicturePath: permanentName,
			IsPrimary:   primary,
			CreatedAt:   helper.TimeNowWIB(),
		})
	}

	return images, errs
}

func (s *productService) Create(req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	product := entity.Product{
		ProductID:   uuid.New().String(),
		ProductCode: req.ProductCode,
		ProductName: req.ProductName,
		ProductSlug: slug.Make(req.ProductName),
		WeightGram:  req.WeightGram,
		TypeID:      req.TypeID,
		Description: req.Description,
		Status:      1,
		CreatedAt:   helper.TimeNowWIB(),
		UpdatedAt:   helper.TimeNowWIB(),
	}

	if err := s.repo.CreateProduct(&product); err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	go func() {

		images, errs := s.buildImages(product.ProductID, req.Images)
		if len(errs) > 0 {
			//successImages := len(images)

			// send only successImages of len(req.Images) saved
		}

		if len(images) > 0 {
			s.repo.CreateProductImages(images)
		}

	}()

	return s.GetByID(product.ProductID)
}

func (s *productService) Update(productID string, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := s.repo.GetProductByID(productID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "product not found"}
	}

	product.ProductName = req.ProductName
	product.ProductSlug = slug.Make(req.ProductName)
	product.WeightGram = req.WeightGram
	product.TypeID = req.TypeID
	product.Description = req.Description
	product.UpdatedAt = helper.TimeNowWIB()

	if err := s.repo.Update(product); err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	go func() {
		if len(req.Images) > 0 {
			oldImages, _ := s.repo.DeleteImagesByProductID(productID)
			for _, img := range oldImages {
				s.minio.RemoveObject(context.Background(), s.bucket, img.PicturePath, minio.RemoveObjectOptions{})
			}

			images, errs := s.buildImages(product.ProductID, req.Images)
			if len(errs) > 0 {
				//successImages := len(images)
				// send only successImages of len(req.Images) saved
			}

			if len(images) > 0 {
				s.repo.CreateProductImages(images)
			}

			//jangan lupa sini rewrite
		}
	}()

	return s.GetByID(productID)
}

func (s *productService) GetByID(productID string) (*dto.ProductResponse, error) {
	product, err := s.repo.GetProductByID(productID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "product not found"}
	}

	images, err := s.repo.GetProductImageByProductID(productID)
	if err != nil {
		images = []entity.ProductImage{}
	}

	var imageResponses []dto.ProductImageResponse

	for _, img := range images {
		imageResponses = append(imageResponses, dto.ProductImageResponse{
			ImageID:     img.ImageID,
			PicturePath: img.PicturePath,
			IsPrimary:   img.IsPrimary,
		})
	}

	return &dto.ProductResponse{
		ProductID:   product.ProductID,
		ProductCode: product.ProductCode,
		ProductName: product.ProductName,
		ProductSlug: product.ProductSlug,
		WeightGram:  product.WeightGram,
		TypeID:      product.TypeID,
		Description: product.Description,
		Status:      product.Status,
		Images:      imageResponses,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil

}

func (s *productService) GetAll() ([]dto.ProductResponse, error) {
	products, err := s.repo.GetAll()
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	var responses []dto.ProductResponse
	for _, p := range products {
		responses = append(responses, dto.ProductResponse{
			ProductID:   p.ProductID,
			ProductCode: p.ProductCode,
			ProductName: p.ProductName,
			ProductSlug: p.ProductSlug,
			WeightGram:  p.WeightGram,
			TypeID:      p.TypeID,
			Description: p.Description,
			Status:      p.Status,
			Images:      []dto.ProductImageResponse{},
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	return responses, nil
}

func (s *productService) Delete(productID string) error {
	_, err := s.repo.GetProductByID(productID)
	if err != nil {
		return &errorhandler.NotFoundError{Message: "product not found"}
	}

	return s.repo.Delete(productID)
}
