package service

import (
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"backend/repository"
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type StoreConfigService interface {
	Upsert(req *dto.StoreConfigRequest) error
	GetConfig() (*dto.StoreConfigResponse, error)
}

type storeConfigService struct {
	repository        repository.StoreConfigRepository
	rajaOngkirService RajaOngkirService
	redis             *redis.Client
}

func NewStoreConfigService(
	repository repository.StoreConfigRepository,
	rajaOngkirService RajaOngkirService,
	redis *redis.Client,
) *storeConfigService {
	return &storeConfigService{
		repository:        repository,
		rajaOngkirService: rajaOngkirService,
		redis:             redis,
	}
}

func (s *storeConfigService) validateAndResolveLocation(provinceID, cityID, districtID, subDistrictID string) (
	provIDStr, provName, cityIDStr, cityName, distIDStr, distName, subDistIDStr, subDistName, zipCode string, err error,
) {
	provIDStr, provName, err = s.rajaOngkirService.FindProvinceByID(provinceID)
	if err != nil {
		return
	}

	cityIDStr, cityName, err = s.rajaOngkirService.FindCityByID(provinceID, cityID)
	if err != nil {
		return
	}

	distIDStr, distName, err = s.rajaOngkirService.FindDistrictByID(cityID, districtID)
	if err != nil {
		return
	}

	subDistIDStr, subDistName, zipCode, err = s.rajaOngkirService.FindSubDistrictByID(districtID, subDistrictID)
	return
}

func (s *storeConfigService) Upsert(req *dto.StoreConfigRequest) error {
	if req.ShopName == "" {
		return &errorhandler.BadRequestError{Message: "Nama toko tidak boleh kosong"}
	}
	if req.Phone == "" {
		return &errorhandler.BadRequestError{Message: "Nomor telepon tidak boleh kosong"}
	}
	if req.ProvinceID == "" {
		return &errorhandler.BadRequestError{Message: "Province ID tidak boleh kosong"}
	}
	if req.CityID == "" {
		return &errorhandler.BadRequestError{Message: "City ID tidak boleh kosong"}
	}
	if req.DistrictID == "" {
		return &errorhandler.BadRequestError{Message: "District ID tidak boleh kosong"}
	}
	if req.SubDistrictID == "" {
		return &errorhandler.BadRequestError{Message: "SubDistrict ID tidak boleh kosong"}
	}

	provIDStr, provName, cityIDStr, cityName, distIDStr, distName, subDistIDStr, subDistName, zipCode, err := s.validateAndResolveLocation(
		req.ProvinceID, req.CityID, req.DistrictID, req.SubDistrictID,
	)
	if err != nil {
		return err
	}

	existing, err := s.repository.GetConfig()
	if err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	if existing == nil {
		config := entity.StoreConfig{
			ConfigID:          uuid.New().String(),
			ShopName:          req.ShopName,
			Phone:             req.Phone,
			ProvinceID:        provIDStr,
			ProvinceName:      provName,
			CityID:            cityIDStr,
			CityName:          cityName,
			DistrictID:        distIDStr,
			DistrictName:      distName,
			SubDistrictID:     subDistIDStr,
			SubDistrictName:   subDistName,
			ZipCode:           zipCode,
			AdditionalAddress: req.AdditionalAddress,
		}
		return s.repository.CreateConfig(&config)
	}

	existing.ShopName = req.ShopName
	existing.Phone = req.Phone
	existing.ProvinceID = provIDStr
	existing.ProvinceName = provName
	existing.CityID = cityIDStr
	existing.CityName = cityName
	existing.DistrictID = distIDStr
	existing.DistrictName = distName
	existing.SubDistrictID = subDistIDStr
	existing.SubDistrictName = subDistName
	existing.ZipCode = zipCode
	existing.AdditionalAddress = req.AdditionalAddress

	ctx := context.Background()
	s.redis.Del(ctx, "config")

	return s.repository.UpdateConfig(existing)
}

func (s *storeConfigService) GetConfig() (*dto.StoreConfigResponse, error) {
	ctx := context.Background()
	cacheKey := "config"

	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var response dto.StoreConfigResponse
		json.Unmarshal([]byte(cached), &response)
		return &response, nil
	}

	if err != redis.Nil {
		return nil, &errorhandler.InternalServerError{Message: "Cache sedang bermasalah, coba lagi"}
	}

	config, err := s.repository.GetConfig()
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	if config == nil {
		return nil, &errorhandler.NotFoundError{Message: "Store config belum diatur"}
	}

	response := &dto.StoreConfigResponse{
		ConfigID:          config.ConfigID,
		ShopName:          config.ShopName,
		Phone:             config.Phone,
		ProvinceID:        config.ProvinceID,
		ProvinceName:      config.ProvinceName,
		CityID:            config.CityID,
		CityName:          config.CityName,
		DistrictID:        config.DistrictID,
		DistrictName:      config.DistrictName,
		SubDistrictID:     config.SubDistrictID,
		SubDistrictName:   config.SubDistrictName,
		ZipCode:           config.ZipCode,
		AdditionalAddress: config.AdditionalAddress,
	}

	jsonData, _ := json.Marshal(response)
	s.redis.Set(ctx, cacheKey, jsonData, 24*time.Hour)

	return response, nil
}
