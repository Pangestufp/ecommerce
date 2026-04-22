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

type ProductPriceService interface {
	Create(req *dto.CreateProductPriceRequest, userID string) (*dto.ProductPriceResponse, error)
	GetAllByProductID(productID string, cursor *dto.Paginate, limit int) ([]dto.ProductPriceResponse, *dto.Paginate, error)
}

type productPriceService struct {
	repository        repository.ProductPriceRepository
	productRepository repository.ProductRepository
	userRepository    repository.UserRepository
	redis             *redis.Client
}

func NewProductPriceService(repository repository.ProductPriceRepository, productRepository repository.ProductRepository, userRepository repository.UserRepository, redis *redis.Client) *productPriceService {
	return &productPriceService{
		repository:        repository,
		productRepository: productRepository,
		userRepository:    userRepository,
		redis:             redis,
	}
}

func (s *productPriceService) Create(req *dto.CreateProductPriceRequest, userID string) (*dto.ProductPriceResponse, error) {
	user, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "User Invalid"}
	}

	_, err = s.productRepository.GetProductByID(req.ProductID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "Product Not Found"}
	}

	price := entity.ProductPrice{
		PriceID:      uuid.New().String(),
		ProductID:    req.ProductID,
		ProductPrice: req.ProductPrice,
		CreatedAt:    helper.TimeNowWIB(),
		CreatedBy:    userID,
		CreatedName:  user.Name,
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

func (s *productPriceService) GetAllByProductID(productID string, cursor *dto.Paginate, limit int) ([]dto.ProductPriceResponse, *dto.Paginate, error) {
	_, err := s.productRepository.GetProductByID(productID)
	if err != nil {
		return nil, nil, &errorhandler.NotFoundError{Message: "Product Not Found"}
	}

	prices, err := s.repository.GetAllByProductID(productID, cursor, limit)
	if err != nil {
		return nil, nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	var paginate *dto.Paginate

	if len(prices) > 0 {
		isNext := cursor == nil || cursor.Direction == nil || *cursor.Direction == "next"
		isPrev := cursor != nil && cursor.Direction != nil && *cursor.Direction == "prev"

		hasNext := "false"
		hasPrev := "false"

		if isNext {
			if len(prices) > limit {
				hasNext = "true"
				prices = prices[:limit]
			}
			if cursor != nil && cursor.LastID != nil {
				hasPrev = "true"
			}
		} else if isPrev {
			if len(prices) > limit {
				hasPrev = "true"
				prices = prices[1:]
			}
			hasNext = "true"
		}

		direction := "next"
		if isPrev {
			direction = "prev"
		}

		first := prices[0]
		last := prices[len(prices)-1]
		paginate = &dto.Paginate{
			FirstID:        &first.PriceID,
			FirstCreatedAt: &first.CreatedAt,
			LastID:         &last.PriceID,
			LastCreatedAt:  &last.CreatedAt,
			HasNext:        &hasNext,
			HasPrev:        &hasPrev,
			Direction:      &direction,
		}
	}

	responses := make([]dto.ProductPriceResponse, 0, len(prices))
	for _, p := range prices {
		responses = append(responses, dto.ProductPriceResponse{
			PriceID:      p.PriceID,
			ProductID:    p.ProductID,
			ProductPrice: p.ProductPrice,
			CreatedAt:    p.CreatedAt,
			CreatedBy:    p.CreatedBy,
			CreatedName:  p.CreatedName,
		})
	}

	return responses, paginate, nil
}
