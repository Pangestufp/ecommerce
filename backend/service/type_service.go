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
	CreateType(req *dto.TypeRequest, userID string) (*dto.TypeResponse, error)
	UpdateType(typeID string, req *dto.TypeRequest, userID string) (*dto.TypeResponse, error)
	DeleteType(typeID string, userID string) error
	GetAllTypePaginate(cursor *dto.Paginate, search string, limit int) ([]dto.TypeResponse, *dto.Paginate, error)
	GetAllType() ([]dto.TypeResponse, error)
	GetTypeByID(typeID string) (*dto.TypeResponse, error)
}

type typeService struct {
	repository repository.TypeRepository
	userRepository repository.UserRepository // Taambahan
	logRepository  repository.LogRepository  // tambahan
	redis      *redis.Client
}

func NewTypeService(repository repository.TypeRepository, redis *redis.Client,userRepository repository.UserRepository,
	logRepository repository.LogRepository ) *typeService {
	return &typeService{repository: repository, redis: redis, userRepository: userRepository,logRepository: logRepository}
}

func (s *typeService) CreateType(req *dto.TypeRequest, userID string) (*dto.TypeResponse, error) {

	user, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "User Invalid"}
	}
	if req.TypeCode == "" {
		return nil, &errorhandler.BadRequestError{Message: "Type Code kosong"}
	}

	if req.TypeName == "" {
		return nil, &errorhandler.BadRequestError{Message: "Type Name kosong"}
	}

	if req.TypeDesc == "" {
		return nil, &errorhandler.BadRequestError{Message: "Type Description kosong"}
	}
	t := entity.Type{
		TypeID:    uuid.New().String(),
		TypeCode:  helper.UpperAndTrim(req.TypeCode),
		TypeName:  req.TypeName,
		TypeDesc:  req.TypeDesc,
		CreatedAt: helper.TimeNowWIB(),
		UpdatedAt: helper.TimeNowWIB(),
		Status:    1,
	}

	data, err := s.repository.GetTypeByTypeCode(helper.UpperAndTrim(req.TypeCode))
	if err == nil && data != nil {
		return nil, &errorhandler.ForbiddenError{Message: "Type Code Telah digunakan"}
	}

	if err := s.repository.CreateType(&t); err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	 s.logRepository.Create(&entity.Log{
		LogID:         uuid.New().String(),
		ReferenceType: "TYPE",
		ReferenceID:   t.TypeID,
		ReferenceName: t.TypeName,
		Note:          fmt.Sprintf("Membuat kategori/tipe baru: %s", t.TypeName),
		CreatedAt:     helper.TimeNowWIB(),
		CreatedBy:     userID,
		CreatedName:   user.Name,
		SourceID:      t.TypeID,
		SourceName:    "CREATE_TYPE",
		SourceType:    "TYPE",
	})


	response := dto.TypeResponse{
		TypeID:   t.TypeID,
		TypeCode: t.TypeCode,
		TypeName: t.TypeName,
		TypeDesc: t.TypeDesc,
	}

	ctx := context.Background()

	cacheKey := "type:all"
	s.redis.Del(ctx, cacheKey)

	return &response, nil
}

func (s *typeService) UpdateType(typeID string, req *dto.TypeRequest, userID string) (*dto.TypeResponse, error) {
	user, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "User Invalid"}
	}
	if req.TypeCode == "" {
		return nil, &errorhandler.BadRequestError{Message: "Type Code kosong"}
	}

	if req.TypeName == "" {
		return nil, &errorhandler.BadRequestError{Message: "Type Name kosong"}
	}

	if req.TypeDesc == "" {
		return nil, &errorhandler.BadRequestError{Message: "Type Description kosong"}
	}

	t, err := s.repository.GetTypeByID(typeID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "Type Not Found"}
	}
	oldName := t.TypeName// buat catete log nya 

	if helper.UpperAndTrim(req.TypeCode) != helper.UpperAndTrim(t.TypeCode) {
		existing, _ := s.repository.GetTypeByTypeCode(helper.UpperAndTrim(req.TypeCode))
		if existing != nil {
			return nil, &errorhandler.ForbiddenError{Message: "Type Code Telah digunakan"}
		}
	}

	t.TypeCode = helper.UpperAndTrim(req.TypeCode)
	t.TypeName = req.TypeName
	t.TypeDesc = req.TypeDesc
	t.UpdatedAt = helper.TimeNowWIB()

	if err := s.repository.UpdateType(t); err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	s.logRepository.Create(&entity.Log{
		LogID:         uuid.New().String(),
		ReferenceType: "TYPE",
		ReferenceID:   t.TypeID,
		ReferenceName: t.TypeName,
		Note:          fmt.Sprintf("Mengubah tipe dari %s menjadi %s", oldName, t.TypeName),
		CreatedAt:     helper.TimeNowWIB(),
		CreatedBy:     userID,
		CreatedName:   user.Name,
		SourceID:      t.TypeID,
		SourceName:    "UPDATE_TYPE",
		SourceType:    "TYPE",
	})

	ctx := context.Background()
	cacheKey := fmt.Sprintf("type:%s", typeID)
	s.redis.Del(ctx, cacheKey)

	cacheKey = "type:all"
	s.redis.Del(ctx, cacheKey)

	return s.GetTypeByID(typeID)

}

func (s *typeService) DeleteType(typeID string, userID string) error {
	t, err := s.repository.GetTypeByID(typeID)//

	if err != nil {
		return &errorhandler.NotFoundError{Message: "type not found"}
	}
	user, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return &errorhandler.NotFoundError{Message: "User Invalid"}
	}

	s.logRepository.Create(&entity.Log{
		LogID:         uuid.New().String(),
		ReferenceType: "TYPE",
		ReferenceID:   t.TypeID,
		ReferenceName: t.TypeName,
		Note:          fmt.Sprintf("Menghapus kategori/tipe: %s", t.TypeName),
		CreatedAt:     helper.TimeNowWIB(),
		CreatedBy:     userID,
		CreatedName:   user.Name,
		SourceID:      t.TypeID,
		SourceName:    "DELETE_TYPE",
		SourceType:    "TYPE",
	})

	ctx := context.Background()
	cacheKey := fmt.Sprintf("type:%s", typeID)
	s.redis.Del(ctx, cacheKey)
	cacheKey = "type:all"
	s.redis.Del(ctx, cacheKey)

	return s.repository.DeleteType(typeID)
}

func (s *typeService) GetAllTypePaginate(cursor *dto.Paginate, search string, limit int) ([]dto.TypeResponse, *dto.Paginate, error) {

	var err error

	types, err := s.repository.GetAllTypePaginate(cursor, search, limit)
	if err != nil {
		return nil, nil, &errorhandler.NotFoundError{Message: "Data Not Found"}
	}

	var paginate *dto.Paginate

	if len(types) > 0 {
		isNext := cursor == nil || cursor.Direction == nil || *cursor.Direction == "next"
		isPrev := cursor != nil && cursor.Direction != nil && *cursor.Direction == "prev"

		hasNext := "false"
		hasPrev := "false"

		if isNext {
			if len(types) > limit {
				hasNext = "true"
				types = types[:limit]
			}
			if cursor != nil && cursor.LastID != nil {
				hasPrev = "true"
			}
		} else if isPrev {
			if len(types) > limit {
				hasPrev = "true"
				types = types[1:]
			}
			hasNext = "true"
		}

		direction := "next"
		if isPrev {
			direction = "prev"
		}

		first := types[0]
		last := types[len(types)-1]
		paginate = &dto.Paginate{
			FirstID:        &first.TypeID,
			FirstCreatedAt: &first.CreatedAt,
			LastID:         &last.TypeID,
			LastCreatedAt:  &last.CreatedAt,
			HasNext:        &hasNext,
			HasPrev:        &hasPrev,
			Direction:      &direction,
		}
	}

	responses := make([]dto.TypeResponse, 0, len(types))
	for _, t := range types {
		responses = append(responses, dto.TypeResponse{
			TypeID:   t.TypeID,
			TypeCode: t.TypeCode,
			TypeName: t.TypeName,
			TypeDesc: t.TypeDesc,
		})
	}

	return responses, paginate, nil
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
		return nil, &errorhandler.NotFoundError{Message: "Type Not Found"}
	}

	responses := make([]dto.TypeResponse, 0, len(types))
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
