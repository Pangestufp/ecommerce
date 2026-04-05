package service

import (
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"backend/repository"

	"github.com/google/uuid"
)

type StoreConfigService interface {
	Upsert(req *dto.StoreConfigRequest) error
	GetConfig() (*dto.StoreConfigResponse, error)
}

type storeConfigService struct {
	repository repository.StoreConfigRepository
}

func NewStoreConfigService(repository repository.StoreConfigRepository) *storeConfigService {
	return &storeConfigService{repository: repository}
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

	return s.repository.UpdateConfig(existing)
}

func (s *storeConfigService) GetConfig() (*dto.StoreConfigResponse, error) {
	config, err := s.repository.GetConfig()
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	if config == nil {
		return nil, &errorhandler.NotFoundError{Message: "store config belum diatur"}
	}

	return &dto.StoreConfigResponse{
		ConfigID: config.ConfigID,
		Origin:   config.Origin,
		Address:  config.Address,
		ShopName: config.ShopName,
		CityID:   config.CityID,
	}, nil
}
