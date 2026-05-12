package service

import (
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"backend/helper"
	"backend/repository"

	"github.com/google/uuid"
)

type AddressService interface {
	CreateAddress(req *dto.CreateAddressRequest, userID string) (*dto.AddressResponse, error)
	UpdateAddress(addressID string, req *dto.UpdateAddressRequest, userID string) (*dto.AddressResponse, error)
	DeleteAddress(addressID string, userID string) error
	GetAddressByUserID(userID string) ([]dto.AddressResponse, error)
}

type addressService struct {
	repository        repository.AddressRepository
	userRepository    repository.UserRepository
	rajaOngkirService RajaOngkirService
}

func NewAddressService(
	repository repository.AddressRepository,
	userRepository repository.UserRepository,
	rajaOngkirService RajaOngkirService,
) *addressService {
	return &addressService{
		repository:        repository,
		userRepository:    userRepository,
		rajaOngkirService: rajaOngkirService,
	}
}

func (s *addressService) validateAndResolveLocation(provinceID, cityID, districtID, subDistrictID string) (
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

func (s *addressService) CreateAddress(req *dto.CreateAddressRequest, userID string) (*dto.AddressResponse, error) {
	_, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "User tidak ditemukan"}
	}

	if req.Label == "" {
		return nil, &errorhandler.BadRequestError{Message: "Label tidak boleh kosong"}
	}
	if req.RecipientName == "" {
		return nil, &errorhandler.BadRequestError{Message: "Nama penerima tidak boleh kosong"}
	}
	if req.Phone == "" {
		return nil, &errorhandler.BadRequestError{Message: "Nomor telepon tidak boleh kosong"}
	}
	if req.ProvinceID == "" {
		return nil, &errorhandler.BadRequestError{Message: "Province ID tidak boleh kosong"}
	}
	if req.CityID == "" {
		return nil, &errorhandler.BadRequestError{Message: "City ID tidak boleh kosong"}
	}
	if req.DistrictID == "" {
		return nil, &errorhandler.BadRequestError{Message: "District ID tidak boleh kosong"}
	}
	if req.SubDistrictID == "" {
		return nil, &errorhandler.BadRequestError{Message: "SubDistrict ID tidak boleh kosong"}
	}

	provIDStr, provName, cityIDStr, cityName, distIDStr, distName, subDistIDStr, subDistName, zipCode, err := s.validateAndResolveLocation(
		req.ProvinceID, req.CityID, req.DistrictID, req.SubDistrictID,
	)
	if err != nil {
		return nil, err
	}

	count, err := s.repository.CountAddressByUserID(userID)
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}
	if count >= 3 {
		return nil, &errorhandler.ForbiddenError{Message: "Maksimal 3 alamat per pengguna"}
	}

	isPrimary := 0
	if req.IsPrimary == 1 || count == 0 {
		if err := s.repository.UnsetPrimaryByUserID(userID); err != nil {
			return nil, &errorhandler.InternalServerError{Message: err.Error()}
		}
		isPrimary = 1
	}

	a := entity.UserAddress{
		AddressID:         uuid.New().String(),
		UserID:            userID,
		Label:             req.Label,
		RecipientName:     req.RecipientName,
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
		IsPrimary:         isPrimary,
		CreatedAt:         helper.TimeNowWIB(),
		UpdatedAt:         helper.TimeNowWIB(),
	}

	if err := s.repository.CreateAddress(&a); err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	return mapToAddressResponse(a), nil
}

func (s *addressService) UpdateAddress(addressID string, req *dto.UpdateAddressRequest, userID string) (*dto.AddressResponse, error) {
	a, err := s.repository.GetAddressByID(addressID)
	if err != nil {
		return nil, err
	}

	if a.UserID != userID {
		return nil, &errorhandler.ForbiddenError{Message: "Tidak memiliki akses ke alamat ini"}
	}

	if req.Label != "" {
		a.Label = req.Label
	}
	if req.RecipientName != "" {
		a.RecipientName = req.RecipientName
	}
	if req.Phone != "" {
		a.Phone = req.Phone
	}

	if req.ProvinceID != "" || req.CityID != "" || req.DistrictID != "" || req.SubDistrictID != "" {
		provinceID := req.ProvinceID
		cityID := req.CityID
		districtID := req.DistrictID
		subDistrictID := req.SubDistrictID

		if provinceID == "" {
			provinceID = a.ProvinceID
		}
		if cityID == "" {
			cityID = a.CityID
		}
		if districtID == "" {
			districtID = a.DistrictID
		}
		if subDistrictID == "" {
			subDistrictID = a.SubDistrictID
		}

		provIDStr, provName, cityIDStr, cityName, distIDStr, distName, subDistIDStr, subDistName, zipCode, err := s.validateAndResolveLocation(
			provinceID, cityID, districtID, subDistrictID,
		)
		if err != nil {
			return nil, err
		}

		a.ProvinceID = provIDStr
		a.ProvinceName = provName
		a.CityID = cityIDStr
		a.CityName = cityName
		a.DistrictID = distIDStr
		a.DistrictName = distName
		a.SubDistrictID = subDistIDStr
		a.SubDistrictName = subDistName
		a.ZipCode = zipCode
	}

	if req.AdditionalAddress != "" {
		a.AdditionalAddress = req.AdditionalAddress
	}

	if req.IsPrimary == 1 && a.IsPrimary != 1 {
		if err := s.repository.UnsetPrimaryByUserID(userID); err != nil {
			return nil, &errorhandler.InternalServerError{Message: err.Error()}
		}
		a.IsPrimary = 1
	}

	a.UpdatedAt = helper.TimeNowWIB()

	if err := s.repository.UpdateAddress(a); err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	return mapToAddressResponse(*a), nil
}

func (s *addressService) DeleteAddress(addressID string, userID string) error {
	a, err := s.repository.GetAddressByID(addressID)
	if err != nil {
		return err
	}

	if a.UserID != userID {
		return &errorhandler.ForbiddenError{Message: "Tidak memiliki akses ke alamat ini"}
	}

	if err := s.repository.DeleteAddress(addressID); err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	if a.IsPrimary == 1 {
		remaining, err := s.repository.GetAddressByUserID(userID)
		if err == nil && len(remaining) > 0 {
			oldest := remaining[len(remaining)-1]
			oldest.IsPrimary = 1
			oldest.UpdatedAt = helper.TimeNowWIB()
			s.repository.UpdateAddress(&oldest)
		}
	}

	return nil
}

func (s *addressService) GetAddressByUserID(userID string) ([]dto.AddressResponse, error) {
	addresses, err := s.repository.GetAddressByUserID(userID)
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	responses := make([]dto.AddressResponse, 0, len(addresses))
	for _, a := range addresses {
		responses = append(responses, *mapToAddressResponse(a))
	}
	return responses, nil
}

func mapToAddressResponse(a entity.UserAddress) *dto.AddressResponse {
	return &dto.AddressResponse{
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
	}
}
