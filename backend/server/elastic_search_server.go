package server

import (
	"backend/config"
	"backend/dto"
	"backend/repository"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type ElasticSearchServer struct {
	DB               *gorm.DB
	ESClient         *elasticsearch.Client
	ProductEventChan chan *dto.ProductEvent
	ProductRepo      repository.ProductRepository
	MinioClient      *minio.Client
	bucket           string
}

var Instance *ElasticSearchServer

func Initialize(db *gorm.DB) {
	Instance = &ElasticSearchServer{
		DB:               db,
		ESClient:         config.ESClient,
		ProductEventChan: make(chan *dto.ProductEvent, 100),
		ProductRepo:      repository.NewProductRepository(db),
		MinioClient:      config.MinioClient,
		bucket:           config.ENV.MinioBucket,
	}
	Instance.startESWriter()
}

func (s *ElasticSearchServer) startESWriter() {
	go func() {
		for event := range s.ProductEventChan {
			enrichedProduct, err := s.ProductRepo.GetProductEnriched(event.ProductID)
			if err != nil {
				log.Printf("Failed to get enriched product %s: %v", event.ProductID, err)
				continue
			}

			if enrichedProduct.Available == 0 {
				if err := s.deleteProductFromES(event.ProductID); err != nil {
					log.Printf("Failed to delete from ES for product %s: %v", event.ProductID, err)
				} else {
					log.Printf("Product %s deleted from ES", event.ProductID)
				}
				continue
			}

			productImages, err := s.ProductRepo.GetProductImageByProductID(event.ProductID)
			if err != nil {
				log.Printf("Failed to get images for product %s: %v", event.ProductID, err)
			}

			images := make([]dto.ProductImageForES, 0, len(productImages))
			for _, image := range productImages {
				images = append(images, dto.ProductImageForES{
					ImageID:     image.ImageID,
					PicturePath: image.PicturePath,
					IsPrimary:   image.IsPrimary,
				})
			}
			enrichedProduct.Images = images

			if err := s.writeProductToES(enrichedProduct); err != nil {
				log.Printf("Failed to write to ES for product %s: %v", event.ProductID, err)
			} else {
				log.Printf("Product %s written to ES (%s)", event.ProductID, event.Type)
			}
		}
	}()
}

func (s *ElasticSearchServer) writeProductToES(product *dto.ProductEnrichedForES) error {
	body, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      "products",
		DocumentID: product.ProductID,
		Body:       strings.NewReader(string(body)),
		Refresh:    "false",
	}

	res, err := req.Do(context.Background(), s.ESClient)
	if err != nil {
		return fmt.Errorf("ES request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("ES error: %s", res.String())
	}

	return nil
}

func (s *ElasticSearchServer) deleteProductFromES(productID string) error {
	req := esapi.DeleteRequest{
		Index:      "products",
		DocumentID: productID,
	}

	res, err := req.Do(context.Background(), s.ESClient)
	if err != nil {
		return fmt.Errorf("ES delete request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("ES delete error: %s", res.String())
	}

	return nil
}

func (s *ElasticSearchServer) SearchProducts(query string) ([]*dto.ProductEnrichedForES, error) {
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"multi_match": map[string]interface{}{
						"query":     query,
						"fields":    []string{"ProductName^3", "ProductCode^2", "ProductSlug^2", "Description", "TypeName"},
						"fuzziness": "AUTO",
					},
				},
				"filter": map[string]interface{}{
					"range": map[string]interface{}{
						"Available": map[string]interface{}{
							"gt": 0,
						},
					},
				},
			},
		},
	}

	body, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	req := esapi.SearchRequest{
		Index: []string{"products"},
		Body:  strings.NewReader(string(body)),
	}

	res, err := req.Do(context.Background(), s.ESClient)
	if err != nil {
		return nil, fmt.Errorf("ES search request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("ES search error: %s", res.String())
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source dto.ProductEnrichedForES `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	products := make([]*dto.ProductEnrichedForES, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		h := hit.Source
		products = append(products, &h)
	}
	s.attachPresignedURLs(products)

	return products, nil
}

func (s *ElasticSearchServer) GetAllProducts() ([]*dto.ProductEnrichedForES, error) {
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				"Available": map[string]interface{}{
					"gt": 0,
				},
			},
		},
	}

	body, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	req := esapi.SearchRequest{
		Index: []string{"products"},
		Body:  strings.NewReader(string(body)),
	}

	res, err := req.Do(context.Background(), s.ESClient)
	if err != nil {
		return nil, fmt.Errorf("ES get all request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("ES get all error: %s", res.String())
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source dto.ProductEnrichedForES `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	products := make([]*dto.ProductEnrichedForES, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		h := hit.Source
		products = append(products, &h)
	}

	s.attachPresignedURLs(products)

	return products, nil
}

func (s *ElasticSearchServer) attachPresignedURLs(products []*dto.ProductEnrichedForES) {
	for _, product := range products {
		validImages := make([]dto.ProductImageForES, 0, len(product.Images))
		for _, image := range product.Images {
			url, err := s.MinioClient.PresignedGetObject(
				context.Background(),
				s.bucket,
				image.PicturePath,
				time.Minute*5,
				nil,
			)
			if err != nil {
				log.Printf("Failed to generate presigned URL for %s: %v", image.PicturePath, err)
				continue
			}
			image.PicturePath = url.String()
			validImages = append(validImages, image)
		}
		product.Images = validImages
	}
}
