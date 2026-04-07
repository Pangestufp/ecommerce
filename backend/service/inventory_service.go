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

type InventoryService interface {
	Create(req *dto.CreateInventoryRequest) (*dto.InventoryResponse, error)
	Update(batchID string, req *dto.UpdateInventoryRequest) (*dto.InventoryResponse, error)
	GetAllByProductID(productID string) ([]dto.InventoryResponse, error)
}

type inventoryService struct {
	repository        repository.InventoryRepository
	productRepository repository.ProductRepository
	redis             *redis.Client
}

func NewInventoryService(repository repository.InventoryRepository, productRepository repository.ProductRepository, redis *redis.Client) *inventoryService {
	return &inventoryService{
		repository:        repository,
		productRepository: productRepository,
		redis:             redis,
	}
}

func (s *inventoryService) Create(req *dto.CreateInventoryRequest) (*dto.InventoryResponse, error) {
	_, err := s.productRepository.GetProductByID(req.ProductID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "Product Not Found"}
	}

	inv := entity.Inventory{
		BatchID:       uuid.New().String(),
		BatchCode:     req.BatchCode,
		ProductID:     req.ProductID,
		CostPrice:     req.CostPrice,
		Stock:         req.Stock,
		ReservedStock: 0,
		CreatedAt:     helper.TimeNowWIB(),
		UpdatedAt:     helper.TimeNowWIB(),
	}

	if err := s.repository.Create(&inv); err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("inventory:%s", inv.ProductID)
	s.redis.Del(ctx, cacheKey)

	return &dto.InventoryResponse{
		BatchID:       inv.BatchID,
		BatchCode:     inv.BatchCode,
		ProductID:     inv.ProductID,
		CostPrice:     inv.CostPrice,
		Stock:         inv.Stock,
		ReservedStock: inv.ReservedStock,
		CreatedAt:     inv.CreatedAt,
		UpdatedAt:     inv.UpdatedAt,
	}, nil
}

func (s *inventoryService) Update(batchID string, req *dto.UpdateInventoryRequest) (*dto.InventoryResponse, error) {
	inv, err := s.repository.GetByID(batchID)
	if err != nil {
		return nil, err
	}

	inv.CostPrice = req.CostPrice
	inv.Stock = req.Stock
	inv.UpdatedAt = helper.TimeNowWIB()

	if err := s.repository.Update(inv); err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("inventory:%s", inv.ProductID)
	s.redis.Del(ctx, cacheKey)

	return &dto.InventoryResponse{
		BatchID:       inv.BatchID,
		BatchCode:     inv.BatchCode,
		ProductID:     inv.ProductID,
		CostPrice:     inv.CostPrice,
		Stock:         inv.Stock,
		ReservedStock: inv.ReservedStock,
		CreatedAt:     inv.CreatedAt,
		UpdatedAt:     inv.UpdatedAt,
	}, nil
}

func (s *inventoryService) GetAllByProductID(productID string) ([]dto.InventoryResponse, error) {

	ctx := context.Background()
	cacheKey := fmt.Sprintf("inventory:%s", productID)

	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var responses []dto.InventoryResponse
		json.Unmarshal([]byte(cached), &responses)
		return responses, nil
	}

	_, err = s.productRepository.GetProductByID(productID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "Product Not Found"}
	}

	inventories, err := s.repository.GetAllByProductID(productID)
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	var responses []dto.InventoryResponse
	for _, inv := range inventories {
		responses = append(responses, dto.InventoryResponse{
			BatchID:       inv.BatchID,
			BatchCode:     inv.BatchCode,
			ProductID:     inv.ProductID,
			CostPrice:     inv.CostPrice,
			Stock:         inv.Stock,
			ReservedStock: inv.ReservedStock,
			CreatedAt:     inv.CreatedAt,
			UpdatedAt:     inv.UpdatedAt,
		})
	}

	jsonData, _ := json.Marshal(responses)
	s.redis.Set(ctx, cacheKey, jsonData, 5*time.Minute)

	return responses, nil
}
