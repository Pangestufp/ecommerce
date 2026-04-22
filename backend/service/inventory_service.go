package service

import (
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"backend/helper"
	"backend/repository"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type InventoryService interface {
	Create(req *dto.CreateInventoryRequest, userID string) (*dto.InventoryResponse, error)
	Update(batchID string, req *dto.UpdateInventoryRequest, userID string) (*dto.InventoryResponse, error)
	GetAllByProductID(productID string, cursor *dto.Paginate, search string, limit int) ([]dto.InventoryResponse, *dto.Paginate, error)
}

type inventoryService struct {
	repository        repository.InventoryRepository
	productRepository repository.ProductRepository
	userRepository    repository.UserRepository
	redis             *redis.Client
}

func NewInventoryService(repository repository.InventoryRepository, productRepository repository.ProductRepository, userRepository repository.UserRepository, redis *redis.Client) *inventoryService {
	return &inventoryService{
		repository:        repository,
		productRepository: productRepository,
		userRepository:    userRepository,
		redis:             redis,
	}
}

func (s *inventoryService) Create(req *dto.CreateInventoryRequest, userID string) (*dto.InventoryResponse, error) {

	user, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "User Invalid"}
	}

	if req.CostPrice < 0 {
		return nil, &errorhandler.BadRequestError{Message: "Modal kurang dari 0"}
	}

	if req.Stock < 0 {
		return nil, &errorhandler.BadRequestError{Message: "jumlah stok minimal adalah 0"}
	}

	product, err := s.productRepository.GetProductByID(req.ProductID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "Product Not Found"}
	}

	seq, yearMonth, err := s.repository.GetNextSeq(product.ProductID)

	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: "Fail to make batch"}
	}

	inv := entity.Inventory{
		BatchID:       uuid.New().String(),
		BatchCode:     fmt.Sprintf("%s-%s-%06d", product.ProductCode, yearMonth, seq),
		ProductID:     req.ProductID,
		CostPrice:     req.CostPrice,
		Stock:         req.Stock,
		ReservedStock: 0,
		CreatedAt:     helper.TimeNowWIB(),
		UpdatedAt:     helper.TimeNowWIB(),
	}

	if err := s.repository.Create(&inv, userID, user.Name); err != nil {
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

func (s *inventoryService) Update(batchID string, req *dto.UpdateInventoryRequest, userID string) (*dto.InventoryResponse, error) {
	user, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "User Invalid"}
	}

	if req.CostPrice < 0 {
		return nil, &errorhandler.BadRequestError{Message: "Modal kurang dari 0"}
	}

	if req.Stock < 0 {
		return nil, &errorhandler.BadRequestError{Message: "jumlah stok minimal adalah 0"}
	}

	inv, err := s.repository.GetByID(batchID)
	if err != nil {
		return nil, err
	}

	inv.CostPrice = req.CostPrice
	inv.Stock = req.Stock
	inv.UpdatedAt = helper.TimeNowWIB()

	if err := s.repository.Update(inv, userID, user.Name); err != nil {
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

func (s *inventoryService) GetAllByProductID(productID string, cursor *dto.Paginate, search string, limit int) ([]dto.InventoryResponse, *dto.Paginate, error) {
	_, err := s.productRepository.GetProductByID(productID)
	if err != nil {
		return nil, nil, &errorhandler.NotFoundError{Message: "Product Not Found"}
	}

	inventories, err := s.repository.GetAllByProductID(productID, cursor, search, limit)
	if err != nil {
		return nil, nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	var paginate *dto.Paginate

	if len(inventories) > 0 {
		isNext := cursor == nil || cursor.Direction == nil || *cursor.Direction == "next"
		isPrev := cursor != nil && cursor.Direction != nil && *cursor.Direction == "prev"

		hasNext := "false"
		hasPrev := "false"

		if isNext {
			if len(inventories) > limit {
				hasNext = "true"
				inventories = inventories[:limit]
			}
			if cursor != nil && cursor.LastID != nil {
				hasPrev = "true"
			}
		} else if isPrev {
			if len(inventories) > limit {
				hasPrev = "true"
				inventories = inventories[1:]
			}
			hasNext = "true"
		}

		direction := "next"
		if isPrev {
			direction = "prev"
		}

		first := inventories[0]
		last := inventories[len(inventories)-1]
		paginate = &dto.Paginate{
			FirstID:        &first.BatchID,
			FirstCreatedAt: &first.CreatedAt,
			LastID:         &last.BatchID,
			LastCreatedAt:  &last.CreatedAt,
			HasNext:        &hasNext,
			HasPrev:        &hasPrev,
			Direction:      &direction,
		}
	}

	responses := make([]dto.InventoryResponse, 0, len(inventories))
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

	return responses, paginate, nil
}
