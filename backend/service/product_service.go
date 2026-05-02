package service

import (
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"backend/helper"
	"backend/repository"
	"backend/server"
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

type ProductService interface {
	GeneratePresignedURLs(req dto.PresignedURLRequest) (*dto.PresignedURLResponse, error)
	Create(req dto.CreateProductRequest) (*dto.ProductResponse, error)
	Update(productID string, req dto.UpdateProductRequest) (*dto.ProductResponse, error)
	GetByID(productID string) (*dto.ProductResponse, error)
	GetAll() ([]dto.ProductResponse, error)
	Delete(productID string) error
	GetAllPaginated(cursor *dto.Paginate, search string, limit int) ([]dto.ProductListRow, *dto.Paginate, error)
	GetProductBySearch(search string, page, limit int) ([]*dto.ProductEnrichedForES, error)
	GetProductEnrichedBySlug(slug string) (*dto.ProductEnrichedForES, error)
}

type productService struct {
	repo   repository.ProductRepository
	repoT  repository.TypeRepository
	repoD  repository.DiscountRepository
	minio  *minio.Client
	redis  *redis.Client
	bucket string
}

func NewProductService(repo repository.ProductRepository, repoT repository.TypeRepository, repoD repository.DiscountRepository, minio *minio.Client, redis *redis.Client, bucket string) *productService {
	return &productService{
		repo:   repo,
		repoT:  repoT,
		repoD:  repoD,
		minio:  minio,
		redis:  redis,
		bucket: bucket,
	}
}

func (s *productService) GeneratePresignedURLs(req dto.PresignedURLRequest) (*dto.PresignedURLResponse, error) {
	var allowedExtensions = map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
		".gif":  true,
	}

	for _, file := range req.Files {
		ext := strings.ToLower(filepath.Ext(file.FileName))
		if !allowedExtensions[ext] {
			return nil, &errorhandler.BadRequestError{
				Message: fmt.Sprintf("tipe file '%s' tidak diizinkan", ext),
			}
		}
	}

	var uploads []dto.UploadItem
	for _, file := range req.Files {
		ext := strings.ToLower(filepath.Ext(file.FileName))
		objectName := "temp/" + uuid.New().String() + ext
		url, err := s.minio.PresignedPutObject(
			context.Background(),
			s.bucket,
			objectName,
			15*time.Minute,
		)
		if err != nil {
			return nil, &errorhandler.InternalServerError{Message: err.Error()}
		}
		uploads = append(uploads, dto.UploadItem{
			UploadURL:  url.String(),
			ObjectName: objectName,
		})
	}

	return &dto.PresignedURLResponse{Uploads: uploads}, nil
}

func (s *productService) moveFromTemp(objectName string) (string, error) {
	obj, err := s.minio.GetObject(context.Background(), s.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return "", &errorhandler.InternalServerError{Message: err.Error()}
	}
	defer obj.Close()

	stat, err := obj.Stat()
	if err != nil {
		return "", &errorhandler.InternalServerError{Message: err.Error()}
	}

	const maxSize = 5 * 1024 * 1024
	if stat.Size > maxSize {
		s.minio.RemoveObject(context.Background(), s.bucket, objectName, minio.RemoveObjectOptions{})
		return "", &errorhandler.BadRequestError{Message: "ukuran file maksimal 5MB"}
	}

	mtype, err := mimetype.DetectReader(obj)
	if err != nil {
		return "", &errorhandler.InternalServerError{Message: err.Error()}
	}

	allowedMimes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
		"image/gif":  true,
	}
	if !allowedMimes[mtype.String()] {
		s.minio.RemoveObject(context.Background(), s.bucket, objectName, minio.RemoveObjectOptions{})
		return "", &errorhandler.BadRequestError{Message: fmt.Sprintf("tipe file '%s' tidak diizinkan", mtype.String())}
	}

	permanentName := "products/" + uuid.New().String() + mtype.Extension()
	src := minio.CopySrcOptions{Bucket: s.bucket, Object: objectName}
	dst := minio.CopyDestOptions{Bucket: s.bucket, Object: permanentName}

	if _, err := s.minio.CopyObject(context.Background(), dst, src); err != nil {
		return "", &errorhandler.InternalServerError{Message: err.Error()}
	}

	s.minio.RemoveObject(context.Background(), s.bucket, objectName, minio.RemoveObjectOptions{})

	return permanentName, nil
}

func (s *productService) buildImages(productID string, objectNames []string) ([]entity.ProductImage, []error) {
	var images []entity.ProductImage
	var errs []error
	primarySet := false

	for _, objectName := range objectNames {
		permanentName, err := s.moveFromTemp(objectName)
		if err != nil {
			errs = append(errs, fmt.Errorf("object %s: %w", objectName, err))
			continue
		}

		primary := 0
		if !primarySet {
			primary = 1
			primarySet = true
		}

		images = append(images, entity.ProductImage{
			ImageID:     uuid.New().String(),
			ProductID:   productID,
			PicturePath: permanentName,
			IsPrimary:   primary,
			CreatedAt:   helper.TimeNowWIB(),
		})
	}

	return images, errs
}

func (s *productService) Create(req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	if req.ProductCode == "" {
		return nil, &errorhandler.BadRequestError{Message: "Product Code kosong"}
	}

	if req.ProductName == "" {
		return nil, &errorhandler.BadRequestError{Message: "Product Name kosong"}
	}

	if req.TypeID == "" {
		return nil, &errorhandler.BadRequestError{Message: "Type ID kosong"}
	}

	if req.Description == "" {
		return nil, &errorhandler.BadRequestError{Message: "Deskripsi kosong"}
	}

	if req.WeightGram < 0 {
		return nil, &errorhandler.BadRequestError{Message: "Weight Gram kosong"}
	}

	data, err := s.repo.GetProductByProductCode(helper.UpperAndTrim(req.ProductCode))

	if err == nil && data != nil {
		return nil, &errorhandler.ForbiddenError{Message: "Product Code Telah digunakan"}
	}

	data, err = s.repo.GetProductByProductSlug(slug.Make(helper.TitleCase(req.ProductName)))
	if err == nil && data != nil {
		return nil, &errorhandler.ForbiddenError{Message: "Product Name Telah digunakan"}
	}

	product := entity.Product{
		ProductID:   uuid.New().String(),
		ProductCode: helper.UpperAndTrim(req.ProductCode),
		ProductName: helper.TitleCase(req.ProductName),
		ProductSlug: slug.Make(helper.TitleCase(req.ProductName)),
		WeightGram:  req.WeightGram,
		TypeID:      req.TypeID,
		Description: req.Description,
		Status:      1,
		CreatedAt:   helper.TimeNowWIB(),
		UpdatedAt:   helper.TimeNowWIB(),
	}

	if err := s.repo.CreateProduct(&product); err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	go func() {

		images, errs := s.buildImages(product.ProductID, req.Images)
		if len(errs) > 0 {
			log.Println("error ", len(errs))
			// send only successImages of len(req.Images) saved
		}

		if len(images) > 0 {
			s.repo.CreateProductImages(images)
		}

		server.Instance.ProductEventChan <- &dto.ProductEvent{
			ProductID: product.ProductID,
			Type:      "create product",
		}

	}()

	return s.GetByID(product.ProductID)
}

func (s *productService) Update(productID string, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {

	if req.ProductCode == "" {
		return nil, &errorhandler.BadRequestError{Message: "Product Code kosong"}
	}

	if req.ProductName == "" {
		return nil, &errorhandler.BadRequestError{Message: "Product Name kosong"}
	}

	if req.TypeID == "" {
		return nil, &errorhandler.BadRequestError{Message: "Type ID kosong"}
	}

	if req.Description == "" {
		return nil, &errorhandler.BadRequestError{Message: "Deskripsi kosong"}
	}

	if req.WeightGram < 0 {
		return nil, &errorhandler.BadRequestError{Message: "Weight Gram kosong"}
	}

	product, err := s.repo.GetProductByID(productID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "product not found"}
	}

	if helper.UpperAndTrim(req.ProductCode) != product.ProductCode {
		data, err := s.repo.GetProductByProductCode(helper.UpperAndTrim(req.ProductCode))
		if err == nil && data != nil {
			return nil, &errorhandler.ForbiddenError{Message: "Product Code Telah digunakan"}
		}
	}

	if helper.TitleCase(req.ProductName) != product.ProductName {
		data, err := s.repo.GetProductByProductSlug(slug.Make(helper.TitleCase(req.ProductName)))
		if err == nil && data != nil {
			return nil, &errorhandler.ForbiddenError{Message: "Product Name Telah digunakan"}
		}
	}

	product.ProductName = helper.TitleCase(req.ProductName)
	product.ProductCode = helper.UpperAndTrim(req.ProductCode)
	product.ProductSlug = slug.Make(helper.TitleCase(req.ProductName))
	product.WeightGram = req.WeightGram
	product.TypeID = req.TypeID
	product.Description = req.Description
	product.UpdatedAt = helper.TimeNowWIB()

	if err := s.repo.Update(product); err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	go func() {
		if len(req.Images) > 0 {
			oldImages, _ := s.repo.DeleteImagesByProductID(productID)
			for _, img := range oldImages {
				s.minio.RemoveObject(context.Background(), s.bucket, img.PicturePath, minio.RemoveObjectOptions{})
			}

			images, errs := s.buildImages(product.ProductID, req.Images)
			if len(errs) > 0 {
				//successImages := len(images)
				// send only successImages of len(req.Images) saved
			}

			if len(images) > 0 {
				s.repo.CreateProductImages(images)
			}

			server.Instance.ProductEventChan <- &dto.ProductEvent{
				ProductID: product.ProductID,
				Type:      "Update product",
			}

			//jangan lupa sini rewrite
		}
	}()

	return s.GetByID(productID)
}

func (s *productService) GetByID(productID string) (*dto.ProductResponse, error) {
	product, err := s.repo.GetProductByID(productID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "product not found"}
	}

	images, err := s.repo.GetProductImageByProductID(productID)
	if err != nil {
		images = []entity.ProductImage{}
	}

	var imageResponses []dto.ProductImageResponse

	ctx := context.Background()

	for _, img := range images {
		cacheKey := fmt.Sprintf("image:%s", img.ImageID)

		cached, err := s.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			imageResponses = append(imageResponses, dto.ProductImageResponse{
				ImageID:     img.ImageID,
				PicturePath: cached,
				IsPrimary:   img.IsPrimary,
			})
			continue
		}

		url, err := s.minio.PresignedGetObject(
			ctx,
			s.bucket,
			img.PicturePath,
			time.Minute*5,
			nil,
		)
		if err != nil {
			log.Printf("Failed to generate presigned URL for %s: %v", img.PicturePath, err)
			continue
		}

		presignedURL := url.String()
		s.redis.Set(ctx, cacheKey, presignedURL, 4*time.Minute)

		imageResponses = append(imageResponses, dto.ProductImageResponse{
			ImageID:     img.ImageID,
			PicturePath: presignedURL,
			IsPrimary:   img.IsPrimary,
		})
	}

	varType, err := s.repoT.GetTypeByID(product.TypeID)
	TypeCode := "error data"
	TypeName := "error data"
	if err == nil {
		TypeCode = varType.TypeCode
		TypeName = varType.TypeName
	}

	return &dto.ProductResponse{
		ProductID:   product.ProductID,
		ProductCode: product.ProductCode,
		ProductName: product.ProductName,
		ProductSlug: product.ProductSlug,
		WeightGram:  product.WeightGram,
		TypeID:      product.TypeID,
		TypeName:    TypeName,
		TypeCode:    TypeCode,
		Description: product.Description,
		Status:      product.Status,
		Images:      imageResponses,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil

}

func (s *productService) GetAll() ([]dto.ProductResponse, error) {

	var responses []dto.ProductResponse

	return responses, nil
}

func (s *productService) Delete(productID string) error {
	_, err := s.repo.GetProductByID(productID)
	if err != nil {
		return &errorhandler.NotFoundError{Message: "product not found"}
	}

	go func() {
		server.Instance.ProductEventChan <- &dto.ProductEvent{
			ProductID: productID,
			Type:      "create product price",
		}
	}()

	return s.repo.Delete(productID)
}

func (s *productService) GetAllPaginated(cursor *dto.Paginate, search string, limit int) ([]dto.ProductListRow, *dto.Paginate, error) {
	products, err := s.repo.GetAllProductsPaginated(cursor, search, limit)
	if err != nil {
		return nil, nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	var paginate *dto.Paginate
	if len(products) > 0 {
		isNext := cursor == nil || cursor.Direction == nil || *cursor.Direction == "next"
		isPrev := cursor != nil && cursor.Direction != nil && *cursor.Direction == "prev"

		hasNext := "false"
		hasPrev := "false"

		if isNext {
			if len(products) > limit {
				hasNext = "true"
				products = products[:limit]
			}
			if cursor != nil && cursor.LastID != nil {
				hasPrev = "true"
			}
		} else if isPrev {
			if len(products) > limit {
				hasPrev = "true"
				products = products[1:]
			}
			hasNext = "true"
		}

		direction := "next"
		if isPrev {
			direction = "prev"
		}

		first := products[0]
		last := products[len(products)-1]
		paginate = &dto.Paginate{
			FirstID:        &first.ProductID,
			FirstCreatedAt: &first.CreatedAt,
			LastID:         &last.ProductID,
			LastCreatedAt:  &last.CreatedAt,
			HasNext:        &hasNext,
			HasPrev:        &hasPrev,
			Direction:      &direction,
		}
	}

	return products, paginate, nil
}

func (s *productService) GetProductBySearch(search string, page, limit int) ([]*dto.ProductEnrichedForES, error) {

	var products []*dto.ProductEnrichedForES
	var err error

	clean := strings.TrimSpace(search)
	from := (page - 1) * limit

	if clean == "" {
		products, err = server.Instance.GetAllProducts(from, limit)
	} else {
		products, err = server.Instance.SearchProducts(search, from, limit)
	}

	return products, err
}

func (s *productService) GetProductEnrichedBySlug(slug string) (*dto.ProductEnrichedForES, error) {
	product, err := s.repo.GetProductByProductSlug(slug)
	if err != nil {
		return nil, &errorhandler.BadRequestError{Message: "data tidak ditemukan"}
	}

	images, err := s.repo.GetProductImageByProductID(product.ProductID)
	if err != nil {
		images = []entity.ProductImage{}
	}

	var imageResponses []dto.ProductImageResponse

	ctx := context.Background()

	for _, img := range images {
		cacheKey := fmt.Sprintf("image:%s", img.ImageID)

		cached, err := s.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			imageResponses = append(imageResponses, dto.ProductImageResponse{
				ImageID:     img.ImageID,
				PicturePath: cached,
				IsPrimary:   img.IsPrimary,
			})
			continue
		}

		url, err := s.minio.PresignedGetObject(
			ctx,
			s.bucket,
			img.PicturePath,
			time.Minute*5,
			nil,
		)
		if err != nil {
			log.Printf("Failed to generate presigned URL for %s: %v", img.PicturePath, err)
			continue
		}

		presignedURL := url.String()
		s.redis.Set(ctx, cacheKey, presignedURL, 4*time.Minute)

		imageResponses = append(imageResponses, dto.ProductImageResponse{
			ImageID:     img.ImageID,
			PicturePath: presignedURL,
			IsPrimary:   img.IsPrimary,
		})
	}

	enrichedproduct, err := s.repo.GetProductEnriched(product.ProductID)

	if err != nil {
		return nil, &errorhandler.BadRequestError{Message: "data tidak ditemukan"}
	}

	enrichedproduct.Images = imageResponses

	return enrichedproduct, err
}
