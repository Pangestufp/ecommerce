package service

import (
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"backend/helper"
	"backend/repository"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

type CheckoutService interface {
	VerifyCheckout(req *dto.CartRequest, userID string) (*dto.VerifyCheckoutResponse, error)
}

type checkoutService struct {
	productRepository  repository.ProductRepository
	discountRepository repository.DiscountRepository
	priceRepository    repository.ProductPriceRepository
	addressRepository  repository.AddressRepository
	minio              *minio.Client
	redis              *redis.Client
	bucket             string
}

func NewCheckoutService(productRepository repository.ProductRepository, discountRepository repository.DiscountRepository, priceRepository repository.ProductPriceRepository, addressRepository repository.AddressRepository, minio *minio.Client, redis *redis.Client, bucket string) *checkoutService {
	return &checkoutService{
		productRepository:  productRepository,
		discountRepository: discountRepository,
		priceRepository:    priceRepository,
		addressRepository:  addressRepository,
		minio:              minio,
		redis:              redis,
		bucket:             bucket,
	}
}

func (s *checkoutService) VerifyCheckout(req *dto.CartRequest, userID string) (*dto.VerifyCheckoutResponse, error) {
	if len(req.ListCart) == 0 {
		return nil, &errorhandler.BadRequestError{Message: "Keranjang tidak valid"}
	}
	if len(req.ListCart) > 20 {
		return nil, &errorhandler.BadRequestError{Message: "Jumlah item keranjang tidak valid"}
	}

	productIDs := make([]string, 0, len(req.ListCart))
	for _, item := range req.ListCart {
		productIDs = append(productIDs, item.ProductID)
	}

	enriched, err := s.productRepository.GetProductsEnrichedBatch(productIDs)
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: "Error validasi item"}
	}

	now := helper.TimeNowWIB()

	allDiscounts, err := s.discountRepository.GetActiveDiscountsByProductIDs(productIDs, now)
	if err != nil {
		allDiscounts = []entity.Discount{}
	}

	allPrices, err := s.priceRepository.GetLatestByProductIDs(productIDs)
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: "Error ambil harga produk"}
	}

	productMap := make(map[string]*dto.ProductEnrichedForES)
	for _, p := range enriched {
		productMap[p.ProductID] = p
	}

	discountMap := make(map[string][]entity.Discount)
	for _, d := range allDiscounts {
		discountMap[d.ProductID] = append(discountMap[d.ProductID], d)
	}

	priceMap := make(map[string]entity.ProductPrice)
	for _, p := range allPrices {
		if _, exists := priceMap[p.ProductID]; !exists {
			priceMap[p.ProductID] = p
		}
	}

	ctx := context.Background()
	var listProduct []dto.ProductCheckoutData

	for _, item := range req.ListCart {
		product, exists := productMap[item.ProductID]
		if !exists || product.Available == 0 {
			continue
		}

		finalQty := item.Qty
		if int64(finalQty) > product.AvailableStock {
			finalQty = int(product.AvailableStock)
		}
		if finalQty == 0 {
			continue
		}

		price, hasPrice := priceMap[product.ProductID]
		if !hasPrice {
			continue
		}

		cacheKey := fmt.Sprintf("image:%s", product.PrimaryImageID)
		presignedURL := ""

		cached, err := s.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			presignedURL = cached
		} else {
			url, err := s.minio.PresignedGetObject(ctx, s.bucket, product.PrimaryImage, time.Minute*5, nil)
			if err != nil {
				log.Printf("Failed to generate presigned URL for %s: %v", product.PrimaryImage, err)
				continue
			}
			presignedURL = url.String()
			s.redis.Set(ctx, cacheKey, presignedURL, 4*time.Minute)
		}

		discountResponses := s.buildDiscountResponses(discountMap[product.ProductID], price, now)

		listProduct = append(listProduct, dto.ProductCheckoutData{
			ProductID:          product.ProductID,
			ProductName:        product.ProductName,
			Image:              presignedURL,
			AvailableStock:     int(product.AvailableStock),
			Qty:                finalQty,
			ProductPrice:       product.ProductPrice,
			ProductPriceFormat: helper.FormatRupiah(product.ProductPrice),
			Discounts:          discountResponses,
		})
	}

	addresses, err := s.addressRepository.GetAddressByUserID(userID)
	if err != nil {
		addresses = []entity.UserAddress{}
	}

	addressResponses := make([]dto.AddressResponse, 0, len(addresses))
	for _, a := range addresses {
		addressResponses = append(addressResponses, dto.AddressResponse{
			AddressID:         a.AddressID,
			UserID:            a.UserID,
			Label:             a.Label,
			RecipientName:     a.RecipientName,
			Phone:             a.Phone,
			ProvinceID:        a.ProvinceID,
			ProvinceName:      a.ProvinceName,
			CityID:            a.CityID,
			CityName:          a.CityName,
			DistrictID:        a.DistrictID,
			DistrictName:      a.DistrictName,
			SubDistrictID:     a.SubDistrictID,
			SubDistrictName:   a.SubDistrictName,
			ZipCode:           a.ZipCode,
			AdditionalAddress: a.AdditionalAddress,
			IsPrimary:         a.IsPrimary,
			CreatedAt:         a.CreatedAt,
			UpdatedAt:         a.UpdatedAt,
		})
	}

	return &dto.VerifyCheckoutResponse{
		ProductPrice: listProduct,
		User_Address: addressResponses,
	}, nil
}

func (s *checkoutService) buildDiscountResponses(discounts []entity.Discount, price entity.ProductPrice, now time.Time) []dto.DiscountResponse {
	responses := make([]dto.DiscountResponse, 0, len(discounts))

	for _, d := range discounts {
		originalPrice := price.ProductPrice

		var discountAmount decimal.Decimal
		if d.DiscountType == helper.Percentage() {
			discountAmount = originalPrice.Mul(d.DiscountValue)
		} else {
			discountAmount = d.DiscountValue
		}

		finalPrice := originalPrice.Sub(discountAmount)

		discountValueFormat := ""
		if d.DiscountType == helper.Percentage() {
			percentage := d.DiscountValue.Mul(decimal.NewFromInt(100))
			discountValueFormat = fmt.Sprintf("%s%%", percentage.StringFixed(0))
		} else {
			discountValueFormat = helper.FormatRupiah(d.DiscountValue)
		}

		statusFormat := "Aktif"
		if now.Before(d.StartAt) {
			statusFormat = "Belum Aktif"
		}

		responses = append(responses, dto.DiscountResponse{
			DiscountID:           d.DiscountID,
			ProductID:            d.ProductID,
			DiscountName:         d.DiscountName,
			DiscountType:         d.DiscountType,
			DiscountValue:        d.DiscountValue,
			DiscountValueFormat:  discountValueFormat,
			DiscountAmountFormat: helper.FormatRupiah(discountAmount),
			FinalValue:           helper.FormatRupiah(finalPrice),
			FinalAmount:          finalPrice,
			StartAt:              d.StartAt,
			ExpiredAt:            d.ExpiredAt,
			StartAtFormat:        helper.FormatTanggalIndo(d.StartAt),
			ExpiredAtFormat:      helper.FormatTanggalIndo(d.ExpiredAt),
			Status:               d.Status,
			StatusFormat:         statusFormat,
			CreatedAt:            d.CreatedAt,
			CreatedBy:            d.CreatedBy,
			CreatedName:          d.CreatedName,
		})
	}

	return responses
}
