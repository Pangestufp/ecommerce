package service

import (
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"backend/helper"
	"backend/repository"
	"backend/server"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

type DiscountService interface {
	Create(req *dto.CreateDiscountRequest, userID string) (*dto.DiscountResponse, error)
	Delete(discountID string, userID string) error
	GetAllByProductID(productID string, cursor *dto.Paginate, search string, limit int) ([]dto.DiscountResponse, *dto.Paginate, error)
	GetDiscountType() []dto.DiscountType
}

type discountService struct {
	repository        repository.DiscountRepository
	productRepository repository.ProductRepository
	userRepository    repository.UserRepository
	priceRepository   repository.ProductPriceRepository
	logRepository     repository.LogRepository
	redis             *redis.Client
}

func NewDiscountService(repository repository.DiscountRepository, productRepository repository.ProductRepository, userRepository repository.UserRepository, priceRepository repository.ProductPriceRepository, logRepository repository.LogRepository, redis *redis.Client) *discountService {
	return &discountService{
		repository:        repository,
		productRepository: productRepository,
		userRepository:    userRepository,
		priceRepository:   priceRepository,
		logRepository:     logRepository,
		redis:             redis,
	}
}

func (s *discountService) Create(req *dto.CreateDiscountRequest, userID string) (*dto.DiscountResponse, error) {
	user, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "User Invalid"}
	}

	product, err := s.productRepository.GetProductByID(req.ProductID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "Product Not Found"}
	}

	price, priceErr := s.priceRepository.GetLatestByProductID(req.ProductID)

	if priceErr != nil {
		return nil, &errorhandler.NotFoundError{Message: "Harga harus sudah disetting"}
	}

	discountValue := decimal.NewFromFloat(req.DiscountValue)
	if req.DiscountType == helper.Amount() && discountValue.GreaterThan(price.ProductPrice) {
		return nil, &errorhandler.BadRequestError{Message: "Diskon melebihi harga jual"}
	}

	if req.DiscountName == "" {
		return nil, &errorhandler.BadRequestError{Message: "Nama diskon kosong"}
	}

	if req.DiscountType != helper.Amount() && req.DiscountType != helper.Percentage() {
		return nil, &errorhandler.BadRequestError{Message: "Tipe diskon tidak valid kosong"}
	}

	if req.DiscountValue <= 0 {
		return nil, &errorhandler.BadRequestError{Message: "Diskon value tidak valid"}
	}

	if req.DiscountType == helper.Percentage() && req.DiscountValue >= 1 {
		return nil, &errorhandler.BadRequestError{Message: "Diskon value tidak valid"}
	}

	startDate, err := time.Parse("2006-01-02", req.StartAt)
	if err != nil {
		return nil, &errorhandler.BadRequestError{Message: "Tanggal tidak valid"}
	}

	endDate, err := time.Parse("2006-01-02", req.ExpiredAt)
	if err != nil {
		return nil, &errorhandler.BadRequestError{Message: "Tanggal tidak valid"}
	}

	startAt := startDate
	expiredAt := endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	timeNow := helper.TimeNowWIB()

	if expiredAt.Before(timeNow) {
		return nil, &errorhandler.BadRequestError{Message: "Tanggal selesai diskon tidak valid"}
	}

	discount := entity.Discount{
		DiscountID:    uuid.New().String(),
		ProductID:     req.ProductID,
		DiscountName:  req.DiscountName,
		DiscountType:  req.DiscountType,
		DiscountValue: decimal.NewFromFloat(req.DiscountValue),
		StartAt:       startAt,
		ExpiredAt:     expiredAt,
		Status:        1,
		CreatedBy:     userID,
		CreatedName:   user.Name,
		CreatedAt:     helper.TimeNowWIB(),
	}

	if err := s.repository.Create(&discount); err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	go func() {

		server.Instance.ProductEventChan <- &dto.ProductEvent{
			ProductID: req.ProductID,
			Type:      "create discount",
		}
	}()

	discountValueFormat := ""
	if discount.DiscountType == helper.Percentage() {
		percentage := discount.DiscountValue.Mul(decimal.NewFromInt(100))
		discountValueFormat = fmt.Sprintf("%s%%", percentage.StringFixed(0))
	} else {
		discountValueFormat = helper.FormatRupiah(discount.DiscountValue)
	}

	discountAmountFormat := "Harga belum diatur"
	finalValue := "Harga belum diatur"

	originalPrice := price.ProductPrice
	var discountAmount decimal.Decimal

	if discount.DiscountValue.LessThan(decimal.NewFromInt(1)) {
		discountAmount = originalPrice.Mul(discount.DiscountValue)
	} else {
		discountAmount = discount.DiscountValue
	}

	finalPrice := originalPrice.Sub(discountAmount)
	discountAmountFormat = helper.FormatRupiah(discountAmount)
	finalValue = helper.FormatRupiah(finalPrice)

	statusFormat := ""
	now := helper.TimeNowWIB()
	if now.Before(discount.StartAt) {
		statusFormat = "Belum Aktif"
	} else {
		statusFormat = "Aktif"
	}

	//create
	ctx := context.Background()
	cacheKey := fmt.Sprintf("ProductDiscount:%s", discount.ProductID)
	s.redis.Del(ctx, cacheKey)

	note := fmt.Sprintf("Membuat diskon '%s' dengan nilai %s", discount.DiscountName, discountValueFormat)
	s.logRepository.Create(&entity.Log{
		LogID:         uuid.New().String(),
		ReferenceType: "PRODUCT",
		ReferenceID:   discount.ProductID,
		ReferenceName: product.ProductName,
		Note:          note,
		CreatedAt:     helper.TimeNowWIB(),
		CreatedBy:     userID,
		CreatedName:   user.Name,
		SourceID:      discount.DiscountID,
		SourceName:    discount.DiscountName,
		SourceType:    "DISCOUNT",
	})

	return &dto.DiscountResponse{
		DiscountID:           discount.DiscountID,
		ProductID:            discount.ProductID,
		DiscountName:         discount.DiscountName,
		DiscountType:         discount.DiscountType,
		DiscountValue:        discount.DiscountValue,
		StartAt:              discount.StartAt,
		ExpiredAt:            discount.ExpiredAt,
		Status:               discount.Status,
		CreatedAt:            discount.CreatedAt,
		DiscountValueFormat:  discountValueFormat,
		DiscountAmountFormat: discountAmountFormat,
		FinalValue:           finalValue,
		StartAtFormat:        helper.FormatTanggalIndo(discount.StartAt),
		ExpiredAtFormat:      helper.FormatTanggalIndo(discount.ExpiredAt),
		StatusFormat:         statusFormat,
		CreatedBy:            discount.CreatedBy,
		CreatedName:          discount.CreatedName,
	}, nil
}

func (s *discountService) Delete(discountID string, userID string) error {
	discount, err := s.repository.GetByID(discountID)
	if err != nil {
		return err
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("ProductDiscount:%s", discount.ProductID)
	s.redis.Del(ctx, cacheKey)

	go func() {
		server.Instance.ProductEventChan <- &dto.ProductEvent{
			ProductID: discount.ProductID,
			Type:      "delete discount",
		}
	}()

	product, err := s.productRepository.GetProductByID(discount.ProductID)
	if err != nil {
		return &errorhandler.NotFoundError{Message: "Product Not Found"}
	}

	user, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return &errorhandler.NotFoundError{Message: "User Invalid"}
	}

	//delete
	note := fmt.Sprintf("Menghapus diskon '%s'", discount.DiscountName)
	s.logRepository.Create(&entity.Log{
		LogID:         uuid.New().String(),
		ReferenceType: "DISCOUNT",
		ReferenceID:   discount.ProductID, // Tetap gunakan Product ID untuk tracking
		ReferenceName: product.ProductName,
		Note:          note,
		CreatedAt:     helper.TimeNowWIB(),
		CreatedBy:     userID,
		CreatedName:   user.Name,
	})

	return s.repository.Delete(discountID)
}

func (s *discountService) GetAllByProductID(productID string, cursor *dto.Paginate, search string, limit int) ([]dto.DiscountResponse, *dto.Paginate, error) {
	_, err := s.productRepository.GetProductByID(productID)
	if err != nil {
		return nil, nil, &errorhandler.NotFoundError{Message: "Product Not Found"}
	}

	discounts, err := s.repository.GetAllByProductID(productID, cursor, search, limit)
	if err != nil {
		return nil, nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	var paginate *dto.Paginate

	if len(discounts) > 0 {
		isNext := cursor == nil || cursor.Direction == nil || *cursor.Direction == "next"
		isPrev := cursor != nil && cursor.Direction != nil && *cursor.Direction == "prev"

		hasNext := "false"
		hasPrev := "false"

		if isNext {
			if len(discounts) > limit {
				hasNext = "true"
				discounts = discounts[:limit]
			}
			if cursor != nil && cursor.LastID != nil {
				hasPrev = "true"
			}
		} else if isPrev {
			if len(discounts) > limit {
				hasPrev = "true"
				discounts = discounts[1:]
			}
			hasNext = "true"
		}

		direction := "next"
		if isPrev {
			direction = "prev"
		}

		first := discounts[0]
		last := discounts[len(discounts)-1]
		paginate = &dto.Paginate{
			FirstID:        &first.DiscountID,
			FirstCreatedAt: &first.CreatedAt,
			LastID:         &last.DiscountID,
			LastCreatedAt:  &last.CreatedAt,
			HasNext:        &hasNext,
			HasPrev:        &hasPrev,
			Direction:      &direction,
		}
	}

	responses := make([]dto.DiscountResponse, 0, len(discounts))

	price, err := s.priceRepository.GetLatestByProductID(productID)

	if err != nil {
		return nil, nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	statusFormat := ""
	now := helper.TimeNowWIB()

	for _, discount := range discounts {

		discountValueFormat := ""
		if discount.DiscountType == helper.Percentage() {
			percentage := discount.DiscountValue.Mul(decimal.NewFromInt(100))
			discountValueFormat = fmt.Sprintf("%s%%", percentage.StringFixed(0))
		} else {
			discountValueFormat = helper.FormatRupiah(discount.DiscountValue)
		}

		discountAmountFormat := "Harga belum diatur"
		finalValue := "Harga belum diatur"

		originalPrice := price.ProductPrice
		var discountAmount decimal.Decimal

		if discount.DiscountValue.LessThan(decimal.NewFromInt(1)) {
			discountAmount = originalPrice.Mul(discount.DiscountValue)
		} else {
			discountAmount = discount.DiscountValue
		}

		finalPrice := originalPrice.Sub(discountAmount)
		discountAmountFormat = helper.FormatRupiah(discountAmount)
		finalValue = helper.FormatRupiah(finalPrice)

		if now.Before(discount.StartAt) {
			statusFormat = "Belum Aktif"
		} else {
			statusFormat = "Aktif"
		}

		responses = append(responses, dto.DiscountResponse{
			DiscountID:           discount.DiscountID,
			ProductID:            productID,
			DiscountName:         discount.DiscountName,
			DiscountType:         discount.DiscountType,
			DiscountValue:        discount.DiscountValue,
			StartAt:              discount.StartAt,
			ExpiredAt:            discount.ExpiredAt,
			Status:               discount.Status,
			CreatedAt:            discount.CreatedAt,
			DiscountValueFormat:  discountValueFormat,
			DiscountAmountFormat: discountAmountFormat,
			FinalValue:           finalValue,
			StartAtFormat:        helper.FormatTanggalIndo(discount.StartAt),
			ExpiredAtFormat:      helper.FormatTanggalIndo(discount.ExpiredAt),
			StatusFormat:         statusFormat,
			CreatedBy:            discount.CreatedBy,
			CreatedName:          discount.CreatedName,
		})

	}
	return responses, paginate, nil
}

func (s *discountService) GetDiscountType() []dto.DiscountType {
	discountTypes := []string{helper.Amount(), helper.Percentage()}

	responses := []dto.DiscountType{}

	for _, discount := range discountTypes {
		responses = append(responses, dto.DiscountType{
			DiscountType: discount,
		})
	}

	return responses
}
