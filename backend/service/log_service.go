package service

import (
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"backend/repository"
)

type LogService interface {
	GetByProductID(productID string, cursor *dto.Paginate, limit int) ([]dto.LogResponse, *dto.Paginate, error)
	GetByReferenceType(cursor *dto.Paginate, limit int) ([]dto.LogResponse, *dto.Paginate, error)
}

type logService struct {
	repository        repository.LogRepository
	productRepository repository.ProductRepository
}

func NewLogService(repository repository.LogRepository, productRepository repository.ProductRepository) *logService {
	return &logService{
		repository:        repository,
		productRepository: productRepository,
	}
}

func (s *logService) GetByProductID(productID string, cursor *dto.Paginate, limit int) ([]dto.LogResponse, *dto.Paginate, error) {
	_, err := s.productRepository.GetProductByID(productID)
	if err != nil {
		return nil, nil, &errorhandler.NotFoundError{Message: "Product Not Found"}
	}

	logs, err := s.repository.GetByReferenceID(productID, cursor, limit)
	if err != nil {
		return nil, nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	return s.formatPaginationResponse(logs, cursor, limit)
}

func (s *logService) GetByReferenceType(cursor *dto.Paginate, limit int) ([]dto.LogResponse, *dto.Paginate, error) {
	logs, err := s.repository.GetByReferenceType(cursor, limit)
	if err != nil {
		return nil, nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	return s.formatPaginationResponse(logs, cursor, limit)
}

func (s *logService) formatPaginationResponse(logs []entity.Log, cursor *dto.Paginate, limit int) ([]dto.LogResponse, *dto.Paginate, error) {
	if limit <= 0 {
		limit = 5
	}

	var paginate *dto.Paginate

	if len(logs) > 0 {
		isNext := cursor == nil || cursor.Direction == nil || *cursor.Direction == "next"
		isPrev := cursor != nil && cursor.Direction != nil && *cursor.Direction == "prev"

		hasNext := "false"
		hasPrev := "false"

		if isNext {

			if len(logs) > limit {
				hasNext = "true"
				logs = logs[:limit]
			}
			if cursor != nil && cursor.LastID != nil {
				hasPrev = "true"
			}
		} else if isPrev {
			if len(logs) > limit {
				hasPrev = "true"
				logs = logs[1:]
			}
			hasNext = "true"
		}

		direction := "next"
		if isPrev {
			direction = "prev"
		}

		first := logs[0]
		last := logs[len(logs)-1]

		paginate = &dto.Paginate{
			FirstID:        &first.LogID,
			FirstCreatedAt: &first.CreatedAt,
			LastID:         &last.LogID,
			LastCreatedAt:  &last.CreatedAt,
			HasNext:        &hasNext,
			HasPrev:        &hasPrev,
			Direction:      &direction,
		}
	}

	responses := make([]dto.LogResponse, 0, len(logs))
	for _, l := range logs {
		responses = append(responses, dto.LogResponse{
			LogID:         l.LogID,
			ReferenceType: l.ReferenceType,
			ReferenceID:   l.ReferenceID,
			ReferenceName: l.ReferenceName,
			Note:          l.Note,
			CreatedAt:     l.CreatedAt,
			CreatedBy:     l.CreatedBy,
			CreatedName:   l.CreatedName,
		})
	}

	return responses, paginate, nil
}
