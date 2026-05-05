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

	// check cache dulu
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var result []any
		json.Unmarshal([]byte(cached), &result)
		return result, nil
	}

	// cache miss, hit API
	result, err := s.fetch(url)
	if err != nil {
		return nil, err
	}

	// simpan ke cache
	encoded, _ := json.Marshal(result)
	s.redis.Set(ctx, cacheKey, encoded, ttl)

	return result, nil
}

func (s *rajaOngkirService) GetProvince() ([]any, error) {
	return s.getWithCache(
		"ongkir:province",
		fmt.Sprintf("%s/destination/province", s.baseURL),
		24*time.Hour,
	)
}

func (s *rajaOngkirService) GetCity(provinceID string) ([]any, error) {
	return s.getWithCache(
		fmt.Sprintf("ongkir:city:%s", provinceID),
		fmt.Sprintf("%s/destination/city/%s", s.baseURL, provinceID),
		24*time.Hour,
	)
}

func (s *rajaOngkirService) GetDistrict(cityID string) ([]any, error) {
	return s.getWithCache(
		fmt.Sprintf("ongkir:district:%s", cityID),
		fmt.Sprintf("%s/destination/district/%s", s.baseURL, cityID),
		24*time.Hour,
	)
}
