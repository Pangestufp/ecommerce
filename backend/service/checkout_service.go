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
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

type CheckoutService interface {
	CreateCheckout(req *dto.CartRequest, userID string) (*dto.CreateCheckoutResponse, error)
	GetCheckout(checkoutID string, userID string) (*dto.VerifyCheckoutResponse, error)
	CalculateShippingFromAddress(req *dto.ShippingRequest, userID string) (*dto.ShippingResponse, error)
}

type checkoutService struct {
	productRepository  repository.ProductRepository
	discountRepository repository.DiscountRepository
	priceRepository    repository.ProductPriceRepository
	addressRepository  repository.AddressRepository
	storeConfigService StoreConfigService
	rajaOngkirService  RajaOngkirService
	minio              *minio.Client
	redis              *redis.Client
	bucket             string
}

func NewCheckoutService(productRepository repository.ProductRepository, discountRepository repository.DiscountRepository, priceRepository repository.ProductPriceRepository, addressRepository repository.AddressRepository, storeConfigService StoreConfigService, rajaOngkirService RajaOngkirService, minio *minio.Client, redis *redis.Client, bucket string) *checkoutService {
	return &checkoutService{
		productRepository:  productRepository,
		discountRepository: discountRepository,
		priceRepository:    priceRepository,
		addressRepository:  addressRepository,
		storeConfigService: storeConfigService,
		rajaOngkirService:  rajaOngkirService,
		minio:              minio,
		redis:              redis,
		bucket:             bucket,
	}
}

func (s *checkoutService) CreateCheckout(req *dto.CartRequest, userID string) (*dto.CreateCheckoutResponse, error) {
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

	var redisItems []dto.CheckoutRedisItem
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

		redisItems = append(redisItems, dto.CheckoutRedisItem{
			ProductID: product.ProductID,
			Qty:       finalQty,
		})
	}

	if len(redisItems) == 0 {
		return nil, &errorhandler.BadRequestError{Message: "Tidak ada item valid untuk checkout"}
	}

	checkoutID := uuid.New().String()
	redisData := dto.CheckoutRedisData{
		CheckoutID: checkoutID,
		UserID:     userID,
		Items:      redisItems,
		CreatedAt:  helper.TimeNowWIB(),
	}

	ctx := context.Background()
	encoded, _ := json.Marshal(redisData)
	cacheKey := fmt.Sprintf("checkout:%s", checkoutID)

	if err := s.redis.Set(ctx, cacheKey, encoded, 30*time.Minute).Err(); err != nil {
		return nil, &errorhandler.InternalServerError{Message: "Gagal menyimpan sesi checkout"}
	}

	return &dto.CreateCheckoutResponse{CheckoutID: checkoutID}, nil
}

func (s *checkoutService) GetCheckout(checkoutID string, userID string) (*dto.VerifyCheckoutResponse, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("checkout:%s", checkoutID)

	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, &errorhandler.NotFoundError{Message: "Sesi checkout tidak ditemukan atau sudah kedaluwarsa"}
	}
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: "Gagal mengambil sesi checkout"}
	}

	var redisData dto.CheckoutRedisData
	if err := json.Unmarshal([]byte(cached), &redisData); err != nil {
		return nil, &errorhandler.InternalServerError{Message: "Data checkout tidak valid"}
	}

	// Validasi kepemilikan
	if redisData.UserID != userID {
		return nil, &errorhandler.BadRequestError{Message: "Akses tidak diizinkan"}
	}

	// Enrich data produk (sama seperti VerifyCheckout lama)
	productIDs := make([]string, 0, len(redisData.Items))
	qtyMap := make(map[string]int)
	for _, item := range redisData.Items {
		productIDs = append(productIDs, item.ProductID)
		qtyMap[item.ProductID] = item.Qty
	}

	enriched, err := s.productRepository.GetProductsEnrichedBatch(productIDs)
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: "Error mengambil data produk"}
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

	var listProduct []dto.ProductCheckoutData
	for _, item := range redisData.Items {
		product, exists := productMap[item.ProductID]
		if !exists || product.Available == 0 {
			continue
		}

		price, hasPrice := priceMap[product.ProductID]
		if !hasPrice {
			continue
		}

		cacheKey := fmt.Sprintf("image:%s", product.PrimaryImageID)
		presignedURL := ""

		imageCached, err := s.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			presignedURL = imageCached
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
			Qty:                qtyMap[product.ProductID],
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

func (s *checkoutService) CalculateShippingFromAddress(req *dto.ShippingRequest, userID string) (*dto.ShippingResponse, error) {
	log.Println("fee bagian 1")

	ctx := context.Background()

	cacheKey := fmt.Sprintf("checkout:%s", req.CheckoutID)
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, &errorhandler.NotFoundError{Message: "Sesi checkout tidak ditemukan atau sudah kedaluwarsa"}
	}
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: "Gagal mengambil sesi checkout"}
	}

	var redisData dto.CheckoutRedisData
	if err := json.Unmarshal([]byte(cached), &redisData); err != nil {
		return nil, &errorhandler.InternalServerError{Message: "Data checkout tidak valid"}
	}
	if redisData.UserID != userID {
		return nil, &errorhandler.BadRequestError{Message: "Akses tidak diizinkan"}
	}

	// Ambil address user
	address, err := s.addressRepository.GetAddressByIDAndUserID(req.AddressID, userID)
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: "Gagal mengambil alamat"}
	}
	if address == nil {
		return nil, &errorhandler.NotFoundError{Message: "Alamat tidak ditemukan"}
	}

	storeConfig, err := s.storeConfigService.GetConfig()
	if err != nil {
		return nil, err
	}

	productIDs := make([]string, 0, len(redisData.Items))
	qtyMap := make(map[string]int)
	for _, item := range redisData.Items {
		productIDs = append(productIDs, item.ProductID)
		qtyMap[item.ProductID] = item.Qty
	}

	totalWeight, err := s.calculateTotalWeight(redisData.Items)
	if err != nil {
		return nil, err
	}

	shippingReq := &dto.ShippingCostRequest{
		Origin:          storeConfig.DistrictID,
		Destination:     address.DistrictID,
		OriginName:      fmt.Sprintf("%s, %s, %s, %s, %s %s", storeConfig.AdditionalAddress, storeConfig.SubDistrictName, storeConfig.DistrictName, storeConfig.CityName, storeConfig.ProvinceName, storeConfig.ZipCode),
		DestinationName: fmt.Sprintf("%s, %s, %s, %s, %s %s", address.AdditionalAddress, address.SubDistrictName, address.DistrictName, address.CityName, address.ProvinceName, address.ZipCode),
		Weight:          totalWeight,
	}

	log.Println("fee bagian 2")

	return s.rajaOngkirService.CalculateShippingCost(shippingReq)
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

func (s *checkoutService) calculateTotalWeight(items []dto.CheckoutRedisItem) (int, error) {
	productIDs := make([]string, 0, len(items))
	for _, item := range items {
		productIDs = append(productIDs, item.ProductID)
	}

	products, err := s.productRepository.GetWeightDataByProductIDs(productIDs)
	if err != nil {
		return 0, &errorhandler.InternalServerError{Message: "Gagal menghitung berat"}
	}

	productMap := make(map[string]entity.Product)
	for _, p := range products {
		productMap[p.ProductID] = p
	}

	totalActualGram := 0
	totalVolumeGram := 0

	for _, item := range items {
		p, exists := productMap[item.ProductID]
		if !exists {
			continue
		}

		// Versi 1: total berat aktual semua item
		totalActualGram += p.WeightGram * item.Qty

		// Versi 2: total berat volume semua item
		volumeKg := float64(p.LengthCm) * float64(p.WidthCm) * float64(p.HeightCm) / 6000.0
		totalVolumeGram += int(volumeKg*1000) * item.Qty
	}

	totalGram := totalActualGram
	if totalVolumeGram > totalActualGram {
		totalGram = totalVolumeGram
	}

	totalKg := totalGram / 1000
	remainderGram := totalGram % 1000

	if remainderGram > 290 {
		totalKg += 1
	}

	if totalKg == 0 {
		totalKg = 1
	}

	return totalKg * 1000, nil
}
