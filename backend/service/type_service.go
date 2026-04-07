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
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type TypeService interface {
	CreateType(req *dto.TypeRequest) (*dto.TypeResponse, error)
	UpdateType(typeID string, req *dto.TypeRequest) (*dto.TypeResponse, error)
	DeleteType(typeID string) error
	GetAllType() ([]dto.TypeResponse, error)
	GetTypeByID(typeID string) (*dto.TypeResponse, error)
}

type typeService struct {
	repository repository.TypeRepository
	redis      *redis.Client
}

func NewTypeService(repository repository.TypeRepository, redis *redis.Client) *typeService {
	return &typeService{repository: repository, redis: redis}
}

func (s *typeService) CreateType(req *dto.TypeRequest) (*dto.TypeResponse, error) {
	t := entity.Type{
		TypeID:    uuid.New().String(),
		TypeCode:  req.TypeCode,
		TypeName:  req.TypeName,
		TypeDesc:  req.TypeDesc,
		CreatedAt: helper.TimeNowWIB(),
		UpdatedAt: helper.TimeNowWIB(),
		Status:    1,
	}

	if err := s.repository.CreateType(&t); err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	response := dto.TypeResponse{
		TypeID:   t.TypeID,
		TypeCode: t.TypeCode,
		TypeName: t.TypeName,
		TypeDesc: t.TypeDesc,
	}

	return &response, nil
}

func (s *typeService) UpdateType(typeID string, req *dto.TypeRequest) (*dto.TypeResponse, error) {
	t, err := s.repository.GetTypeByID(typeID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "Type Not Found"}
	}

	t.TypeCode = req.TypeCode
	t.TypeName = req.TypeName
	t.TypeDesc = req.TypeDesc
	t.UpdatedAt = helper.TimeNowWIB()

	if err := s.repository.UpdateType(t); err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("type:%s", typeID)
	s.redis.Del(ctx, cacheKey)

	return s.GetTypeByID(typeID)

}

func (s *typeService) DeleteType(typeID string) error {
	_, err := s.repository.GetTypeByID(typeID)

	if err != nil {
		return &errorhandler.NotFoundError{Message: "type not found"}
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("type:%s", typeID)
	s.redis.Del(ctx, cacheKey)

	return s.repository.DeleteType(typeID)
}

func (s *typeService) GetAllType() ([]dto.TypeResponse, error) {
	ctx := context.Background()
	cacheKey := "type:all"

	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var responses []dto.TypeResponse
		json.Unmarshal([]byte(cached), &responses)
		return responses, nil
	}

	types, err := s.repository.GetAllType()
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "Data Not Found"}
	}

	var responses []dto.TypeResponse
	for _, t := range types {
		responses = append(responses, dto.TypeResponse{
			TypeID:   t.TypeID,
			TypeCode: t.TypeCode,
			TypeName: t.TypeName,
			TypeDesc: t.TypeDesc,
		})
	}

	jsonData, _ := json.Marshal(responses)
	s.redis.Set(ctx, cacheKey, jsonData, 15*time.Minute)

	return responses, nil
}

func (s *typeService) GetTypeByID(typeID string) (*dto.TypeResponse, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("type:%s", typeID)

	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var response dto.TypeResponse
		json.Unmarshal([]byte(cached), &response)
		return &response, nil
	}

	t, err := s.repository.GetTypeByID(typeID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "Type Not Found"}
	}

	response := dto.TypeResponse{
		TypeID:   t.TypeID,
		TypeCode: t.TypeCode,
		TypeName: t.TypeName,
		TypeDesc: t.TypeDesc,
	}

	jsonData, _ := json.Marshal(response)
	s.redis.Set(ctx, cacheKey, jsonData, 15*time.Minute)

	return &response, nil
}
