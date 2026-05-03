package service

import (
	"backend/dto"
	"backend/errorhandler"
	"backend/repository"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

type CartService interface {
	VerifyCart(req *dto.CartRequest) (*dto.CartResponse, error)
}

type cartService struct {
	productRepository repository.ProductRepository
	minio             *minio.Client
	redis             *redis.Client
	bucket            string
}

func NewCartService(productRepository repository.ProductRepository, minio *minio.Client, redis *redis.Client, bucket string) *cartService {
	return &cartService{
		productRepository: productRepository,
		minio:             minio,
		redis:             redis,
		bucket:            bucket,
	}
}

func (s *cartService) VerifyCart(req *dto.CartRequest) (*dto.CartResponse, error) {
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

	productMap := make(map[string]*dto.ProductEnrichedForES)
	for _, p := range enriched {
		productMap[p.ProductID] = p
	}

	var (
		listProduct []dto.CartVerifiedProduct
		listSave    []dto.NewSave
		totalNow    = decimal.Zero
		notes       string
	)

	isNote := 0

	for _, item := range req.ListCart {
		product, exists := productMap[item.ProductID]

		if !exists {
			continue
		}

		if product.Available == 0 {
			if isNote == 0 {
				notes = "Beberapa produk sudah tidak tersedia"
				isNote = 1
			}
			continue
		}

		finalQty := item.Qty
		if int64(finalQty) > product.AvailableStock {
			finalQty = int(product.AvailableStock)
		}

		if finalQty == 0 {
			if isNote == 0 {
				notes = "Beberapa stok habis dan telah dihapus dari keranjang"
				isNote = 1
			}
			continue
		}

		bestPrice := decimal.NewFromFloat(product.BestPrice)

		ctx := context.Background()
		cacheKey := fmt.Sprintf("image:%s", product.PrimaryImageID)

		cached, err := s.redis.Get(ctx, cacheKey).Result()

		presignedURL := ""
		if err == nil {
			presignedURL = cached
		} else {

			url, err := s.minio.PresignedGetObject(
				ctx,
				s.bucket,
				product.PrimaryImage,
				time.Minute*5,
				nil,
			)
			if err != nil {
				log.Printf("Failed to generate presigned URL for %s: %v", product.PrimaryImage, err)
				continue
			}

			presignedURL := url.String()
			s.redis.Set(ctx, cacheKey, presignedURL, 4*time.Minute)

		}

		listProduct = append(listProduct, dto.CartVerifiedProduct{
			ProductID:      product.ProductID,
			ProductName:    product.ProductName,
			Image:          presignedURL,
			IsAvailable:    product.Available,
			AvailableStock: int(product.AvailableStock),
			Qty:            finalQty,
			BestPrice:      bestPrice,
			PriceFormat:    product.BestPriceFormat,
		})

		listSave = append(listSave, dto.NewSave{
			ID:          product.ProductID,
			ProductName: product.ProductName,
			Qty:         finalQty,
			Image:       presignedURL,
		})

		totalNow = totalNow.Add(bestPrice.Mul(decimal.NewFromInt(int64(finalQty))))
	}

	return &dto.CartResponse{
		ListProduct: listProduct,
		ListSave:    listSave,
		IsNote:      isNote,
		Note:        notes,
		TotalNow:    totalNow,
	}, nil
}
