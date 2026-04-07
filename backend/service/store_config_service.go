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
	repository repository.StoreConfigRepository
	redis      *redis.Client
}

func NewStoreConfigService(repository repository.StoreConfigRepository, redis *redis.Client) *storeConfigService {
	return &storeConfigService{
		repository: repository,
		redis:      redis,
	}
}

func (s *storeConfigService) Upsert(req *dto.StoreConfigRequest) error {
	existing, err := s.repository.GetConfig()
	if err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	if existing == nil {
		config := entity.StoreConfig{
			ConfigID: uuid.New().String(),
			Origin:   req.Origin,
			Address:  req.Address,
			ShopName: req.ShopName,
			CityID:   req.CityID,
		}
		return s.repository.CreateConfig(&config)
	}

	if req.Origin != "" {
		existing.Origin = req.Origin
	}
	if req.Address != "" {
		existing.Address = req.Address
	}
	if req.ShopName != "" {
		existing.ShopName = req.ShopName
	}
	if req.CityID != "" {
		existing.CityID = req.CityID
	}

	ctx := context.Background()
	cacheKey := "config"
	s.redis.Del(ctx, cacheKey)

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

	config, err := s.repository.GetConfig()
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	if config == nil {
		return nil, &errorhandler.NotFoundError{Message: "store config belum diatur"}
	}

	jsonData, _ := json.Marshal(config)
	s.redis.Set(ctx, cacheKey, jsonData, 24*time.Hour)

	return &dto.StoreConfigResponse{
		ConfigID: config.ConfigID,
		Origin:   config.Origin,
		Address:  config.Address,
		ShopName: config.ShopName,
		CityID:   config.CityID,
	}, nil
}
