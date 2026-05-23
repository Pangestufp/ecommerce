package service

import (
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"backend/helper"
	"backend/repository"
	"strings"

	"github.com/google/uuid"
)

type CourierService interface {
	Create(req *dto.CreateCourierRequest) (*dto.CourierResponse, error)
	GetAll() ([]dto.CourierResponse, error)
	Update(id string, req *dto.UpdateCourierRequest) (*dto.CourierResponse, error)
	Toggle(id string) (*dto.CourierResponse, error)
}

type courierService struct {
	repository repository.CourierRepository
}

func NewCourierService(repository repository.CourierRepository) *courierService {
	return &courierService{repository: repository}
}

func (s *courierService) Create(req *dto.CreateCourierRequest) (*dto.CourierResponse, error) {
	code := strings.ToLower(strings.TrimSpace(req.Code))
	name := strings.TrimSpace(req.Name)

	if code == "" {
		return nil, &errorhandler.BadRequestError{Message: "Code tidak boleh kosong"}
	}
	if name == "" {
		return nil, &errorhandler.BadRequestError{Message: "Nama tidak boleh kosong"}
	}

	// cek duplikat code
	existing, _ := s.repository.GetByCode(code)
	if existing != nil {
		return nil, &errorhandler.BadRequestError{Message: "Code kurir sudah terdaftar"}
	}

	courier := entity.Courier{
		ID:        uuid.New().String(),
		Code:      code,
		Name:      name,
		Status:    1,
		CreatedAt: helper.TimeNowWIB(),
		UpdatedAt: helper.TimeNowWIB(),
	}

	if err := s.repository.Create(&courier); err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	return mapToCourierResponse(courier), nil
}

func (s *courierService) GetAll() ([]dto.CourierResponse, error) {
	couriers, err := s.repository.GetAll()
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	responses := make([]dto.CourierResponse, 0, len(couriers))
	for _, c := range couriers {
		responses = append(responses, *mapToCourierResponse(c))
	}
	return responses, nil
}

func (s *courierService) Update(id string, req *dto.UpdateCourierRequest) (*dto.CourierResponse, error) {
	courier, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, &errorhandler.BadRequestError{Message: "Nama tidak boleh kosong"}
	}

	courier.Name = name

	if req.Code != "" {
		code := strings.ToLower(strings.TrimSpace(req.Code))

		existing, _ := s.repository.GetByCode(code)
		if existing != nil && existing.ID != courier.ID {
			return nil, &errorhandler.BadRequestError{Message: "Code kurir sudah terdaftar"}
		}

		courier.Code = code
	}

	courier.UpdatedAt = helper.TimeNowWIB()

	if err := s.repository.Update(courier); err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	return mapToCourierResponse(*courier), nil
}

func (s *courierService) Toggle(id string) (*dto.CourierResponse, error) {
	// pastiin courier ada dulu
	_, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.repository.Toggle(id); err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	// ambil data terbaru setelah toggle
	updated, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return mapToCourierResponse(*updated), nil
}

func mapToCourierResponse(c entity.Courier) *dto.CourierResponse {
	return &dto.CourierResponse{
		ID:              c.ID,
		Code:            c.Code,
		Name:            c.Name,
		Status:          c.Status,
		CreatedAt:       c.CreatedAt,
		UpdatedAt:       c.UpdatedAt,
		CreatedAtFormat: helper.FormatTanggalIndo(c.CreatedAt),
		UpdatedAtFormat: helper.FormatTanggalIndo(c.UpdatedAt),
	}
}
