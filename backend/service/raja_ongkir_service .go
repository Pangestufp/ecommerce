package service

import (
	"backend/errorhandler"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type RajaOngkirService interface {
	GetProvince() ([]any, error)
	GetCity(provinceID string) ([]any, error)
	GetDistrict(cityID string) ([]any, error)
	GetSubDistrict(districtID string) ([]any, error)
	FindProvinceByID(provinceID string) (string, string, error)
	FindCityByID(provinceID string, cityID string) (string, string, error)
	FindDistrictByID(cityID string, districtID string) (string, string, error)
	FindSubDistrictByID(districtID string, subDistrictID string) (string, string, string, error)
}

type rajaOngkirService struct {
	apiKey  string
	baseURL string
	redis   *redis.Client
}

func NewRajaOngkirService(apiKey string, url string, redis *redis.Client) *rajaOngkirService {
	return &rajaOngkirService{
		apiKey:  apiKey,
		baseURL: url,
		redis:   redis,
	}
}

func (s *rajaOngkirService) fetch(url string) ([]any, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}
	req.Header.Set("key", s.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Data []any `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	return result.Data, nil
}

func (s *rajaOngkirService) getWithCache(cacheKey, url string, ttl time.Duration) ([]any, error) {
	ctx := context.Background()

	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var result []any
		json.Unmarshal([]byte(cached), &result)
		return result, nil
	}

	if err != redis.Nil {
		return nil, &errorhandler.InternalServerError{Message: "Cache sedang bermasalah, coba lagi"}
	}

	result, err := s.fetch(url)
	if err != nil {
		return nil, err
	}

	encoded, _ := json.Marshal(result)
	s.redis.Set(ctx, cacheKey, encoded, ttl)

	return result, nil
}

func (s *rajaOngkirService) GetProvince() ([]any, error) {
	return s.getWithCache(
		"ongkir:province",
		fmt.Sprintf("%s/destination/province", s.baseURL),
		7*24*time.Hour,
	)
}

func (s *rajaOngkirService) GetCity(provinceID string) ([]any, error) {
	return s.getWithCache(
		fmt.Sprintf("ongkir:city:%s", provinceID),
		fmt.Sprintf("%s/destination/city/%s", s.baseURL, provinceID),
		7*24*time.Hour,
	)
}

func (s *rajaOngkirService) GetDistrict(cityID string) ([]any, error) {
	return s.getWithCache(
		fmt.Sprintf("ongkir:district:%s", cityID),
		fmt.Sprintf("%s/destination/district/%s", s.baseURL, cityID),
		7*24*time.Hour,
	)
}

func (s *rajaOngkirService) GetSubDistrict(districtID string) ([]any, error) {
	return s.getWithCache(
		fmt.Sprintf("ongkir:sub-district:%s", districtID),
		fmt.Sprintf("%s/destination/sub-district/%s", s.baseURL, districtID),
		7*24*time.Hour,
	)
}

func (s *rajaOngkirService) FindProvinceByID(provinceID string) (string, string, error) {
	provinces, err := s.GetProvince()
	if err != nil {
		return "", "", err
	}

	for _, p := range provinces {
		item, ok := p.(map[string]any)
		if !ok {
			continue
		}
		id, ok := item["id"].(float64)
		if !ok {
			continue
		}
		if fmt.Sprintf("%d", int(id)) == provinceID {
			name, _ := item["name"].(string)
			return provinceID, name, nil
		}
	}

	return "", "", &errorhandler.BadRequestError{Message: "Province ID tidak valid"}
}

func (s *rajaOngkirService) FindCityByID(provinceID string, cityID string) (string, string, error) {
	cities, err := s.GetCity(provinceID)
	if err != nil {
		return "", "", err
	}

	for _, c := range cities {
		item, ok := c.(map[string]any)
		if !ok {
			continue
		}
		id, ok := item["id"].(float64)
		if !ok {
			continue
		}
		if fmt.Sprintf("%d", int(id)) == cityID {
			name, _ := item["name"].(string)
			return cityID, name, nil
		}
	}

	return "", "", &errorhandler.BadRequestError{Message: "City ID tidak valid"}
}

// FindDistrictByID — zip_code tidak diambil disini, dipindah ke subdistrict
func (s *rajaOngkirService) FindDistrictByID(cityID string, districtID string) (string, string, error) {
	districts, err := s.GetDistrict(cityID)
	if err != nil {
		return "", "", err
	}

	for _, d := range districts {
		item, ok := d.(map[string]any)
		if !ok {
			continue
		}
		id, ok := item["id"].(float64)
		if !ok {
			continue
		}
		if fmt.Sprintf("%d", int(id)) == districtID {
			name, _ := item["name"].(string)
			return districtID, name, nil
		}
	}

	return "", "", &errorhandler.BadRequestError{Message: "District ID tidak valid"}
}

func (s *rajaOngkirService) FindSubDistrictByID(districtID string, subDistrictID string) (string, string, string, error) {
	subDistricts, err := s.GetSubDistrict(districtID)
	if err != nil {
		return "", "", "", err
	}

	for _, sd := range subDistricts {
		item, ok := sd.(map[string]any)
		if !ok {
			continue
		}
		id, ok := item["id"].(float64)
		if !ok {
			continue
		}
		if fmt.Sprintf("%d", int(id)) == subDistrictID {
			name, _ := item["name"].(string)
			zipCode, _ := item["zip_code"].(string)
			if zipCode == "0" {
				zipCode = ""
			}
			return subDistrictID, name, zipCode, nil
		}
	}

	return "", "", "", &errorhandler.BadRequestError{Message: "SubDistrict ID tidak valid"}
}
