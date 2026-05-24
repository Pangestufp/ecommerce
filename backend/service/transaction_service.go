package service

import (
	"backend/dto"
	"backend/errorhandler"
	"backend/repository"
)

type TransactionService interface {
	GetAllByBatchID(batchID string, cursor *dto.Paginate, limit int) ([]dto.TransactionResponse, *dto.Paginate, error)
}

type transactionService struct {
	repository    repository.TransactionRepository
	inventoryRepo repository.InventoryRepository
}

func NewTransactionService(repository repository.TransactionRepository, inventoryRepo repository.InventoryRepository) *transactionService {
	return &transactionService{
		repository:    repository,
		inventoryRepo: inventoryRepo,
	}
}

func (s *transactionService) GetAllByBatchID(batchID string, cursor *dto.Paginate, limit int) ([]dto.TransactionResponse, *dto.Paginate, error) {
	_, err := s.inventoryRepo.GetByID(batchID)
	if err != nil {
		return nil, nil, &errorhandler.NotFoundError{Message: "Batch tidak ditemukan"}
	}


	trxs, err := s.repository.GetAllByBatchIDPaginate(batchID, cursor, limit)
	if err != nil {
		return nil, nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	var paginate *dto.Paginate

	if len(trxs) > 0 {
		isNext := cursor == nil || cursor.Direction == nil || *cursor.Direction == "next"
		isPrev := cursor != nil && cursor.Direction != nil && *cursor.Direction == "prev"

		hasNext := "false"
		hasPrev := "false"

		if isNext {
			if len(trxs) > limit {
				hasNext = "true"
				trxs = trxs[:limit] 
			}
			if cursor != nil && cursor.LastID != nil {
				hasPrev = "true"
			}
		} else if isPrev {
			if len(trxs) > limit {
				hasPrev = "true"
				trxs = trxs[1:]
			}
			hasNext = "true"
		}

		direction := "next"
		if isPrev {
			direction = "prev"
		}

		first := trxs[0]
		last := trxs[len(trxs)-1]
		
		paginate = &dto.Paginate{
			FirstID:        &first.TransactionID,
			FirstCreatedAt: &first.CreatedAt,
			LastID:         &last.TransactionID,
			LastCreatedAt:  &last.CreatedAt,
			HasNext:        &hasNext,
			HasPrev:        &hasPrev,
			Direction:      &direction,
		}
	}

	
	responses := make([]dto.TransactionResponse, 0, len(trxs))
	for _, t := range trxs {
		responses = append(responses, dto.TransactionResponse{
			TransactionID: t.TransactionID,
			BatchID:       t.BatchID,
			Type:          t.Type,
			Quantity:      t.Quantity,
			ReferenceType: t.ReferenceType,
			ReferenceID:   t.ReferenceID,
			Note:          t.Note,
			CreatedAt:     t.CreatedAt,
		})
	}

	return responses, paginate, nil
}