package service

import (
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"backend/helper"
	"backend/repository"

	"github.com/google/uuid"
)

type DiscountService interface {
	Create(req *dto.CreateDiscountRequest) (*dto.DiscountResponse, error)
	Delete(discountID string) error
}

type discountService struct {
	repository        repository.DiscountRepository
	productRepository repository.ProductRepository
}

func NewDiscountService(repository repository.DiscountRepository, productRepository repository.ProductRepository) *discountService {
	return &discountService{
		repository:        repository,
		productRepository: productRepository,
	}
}

func (s *discountService) Create(req *dto.CreateDiscountRequest) (*dto.DiscountResponse, error) {
	_, err := s.productRepository.GetProductByID(req.ProductID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "Product Not Found"}
	}

	discount := entity.Discount{
		DiscountID:    uuid.New().String(),
		ProductID:     req.ProductID,
		DiscountName:  req.DiscountName,
		DiscountType:  req.DiscountType,
		DiscountValue: req.DiscountValue,
		StartAt:       req.StartAt,
		ExpiredAt:     req.ExpiredAt,
		Status:        1,
		CreatedAt:     helper.TimeNowWIB(),
	}

	if err := s.repository.Create(&discount); err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	return &dto.DiscountResponse{
		DiscountID:    discount.DiscountID,
		ProductID:     discount.ProductID,
		DiscountName:  discount.DiscountName,
		DiscountType:  discount.DiscountType,
		DiscountValue: discount.DiscountValue,
		StartAt:       discount.StartAt,
		ExpiredAt:     discount.ExpiredAt,
		Status:        discount.Status,
		CreatedAt:     discount.CreatedAt,
	}, nil
}

func (s *discountService) Delete(discountID string) error {
	_, err := s.repository.GetByID(discountID)
	if err != nil {
		return err
	}

	return s.repository.Delete(discountID)
}
