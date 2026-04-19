package service

import (
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"backend/helper"
	"backend/repository"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type ProductPriceService interface {
	Create(req *dto.CreateProductPriceRequest) (*dto.ProductPriceResponse, error)
	GetAllByProductID(productID string) ([]dto.ProductPriceResponse, error)
}

type productPriceService struct {
	repository        repository.ProductPriceRepository
	productRepository repository.ProductRepository
	redis             *redis.Client
}

func NewProductPriceService(repository repository.ProductPriceRepository, productRepository repository.ProductRepository, redis *redis.Client) *productPriceService {
	return &productPriceService{
		repository:        repository,
		productRepository: productRepository,
		redis:             redis,
	}
}

func (s *productPriceService) Create(req *dto.CreateProductPriceRequest) (*dto.ProductPriceResponse, error) {
	_, err := s.productRepository.GetProductByID(req.ProductID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "Product Not Found"}
	}

	price := entity.ProductPrice{
		PriceID:      uuid.New().String(),
		ProductID:    req.ProductID,
		ProductPrice: req.ProductPrice,
		CreatedAt:    helper.TimeNowWIB(),
	}

	if err := s.repository.Create(&price); err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("ProductPrice:%s", price.ProductID)
	s.redis.Del(ctx, cacheKey)

	return &dto.ProductPriceResponse{
		PriceID:      price.PriceID,
		ProductID:    price.ProductID,
		ProductPrice: price.ProductPrice,
		CreatedAt:    price.CreatedAt,
	}, nil
}

func (s *productPriceService) GetAllByProductID(productID string) ([]dto.ProductPriceResponse, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("ProductPrice:%s", productID)

	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var responses []dto.ProductPriceResponse
		json.Unmarshal([]byte(cached), &responses)
		return responses, nil
	}

	_, err = s.productRepository.GetProductByID(productID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "Product Not Found"}
	}

	prices, err := s.repository.GetAllByProductID(productID)
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	var responses []dto.ProductPriceResponse
	for _, p := range prices {
		responses = append(responses, dto.ProductPriceResponse{
			PriceID:      p.PriceID,
			ProductID:    p.ProductID,
			ProductPrice: p.ProductPrice,
			CreatedAt:    p.CreatedAt,
		})
	}

	jsonData, _ := json.Marshal(responses)
	s.redis.Set(ctx, cacheKey, jsonData, 15*time.Minute)

	return responses, nil
}
