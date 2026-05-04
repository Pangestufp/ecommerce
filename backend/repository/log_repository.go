package repository

import (
	"backend/dto"
	"backend/entity"
	"gorm.io/gorm"
)

type LogRepository interface {
	Create(log *entity.Log) error
	GetByReferenceID(refID string, cursor *dto.Paginate, limit int) ([]entity.Log, error)
	GetByReferenceType(cursor *dto.Paginate, limit int) ([]entity.Log, error)
}

type logRepository struct {
	db *gorm.DB
}

func NewLogRepository(db *gorm.DB) *logRepository {
	return &logRepository{db: db}
}


func (r *logRepository) Create(log *entity.Log) error {
	return r.db.Create(log).Error
}


func (r *logRepository) GetByReferenceID(refID string, cursor *dto.Paginate, limit int) ([]entity.Log, error) {
	query := r.db.Model(&entity.Log{}).Where("reference_id = ?", refID)
	return r.applyPagination(query, cursor, limit)
}


func (r *logRepository) GetByReferenceType(cursor *dto.Paginate, limit int) ([]entity.Log, error) {
	query := r.db.Model(&entity.Log{}).Where("reference_type = 'TYPE'")
	return r.applyPagination(query, cursor, limit)
}


func (r *logRepository) applyPagination(query *gorm.DB, cursor *dto.Paginate, limit int) ([]entity.Log, error) {
	if limit <= 0 {
		limit = 5
	}

	var logs []entity.Log

	if cursor != nil {
		if cursor.Direction != nil && *cursor.Direction == "prev" {
			if cursor.FirstID != nil && cursor.FirstCreatedAt != nil {

				query = query.Where("(created_at, log_id) > (?, ?)", cursor.FirstCreatedAt, cursor.FirstID)
				query = query.Order("created_at ASC, log_id ASC")
			}
		} else {
			if cursor.Direction != nil && *cursor.Direction == "next" {
				
				query = query.Where("(created_at, log_id) < (?, ?)", cursor.LastCreatedAt, cursor.LastID)
			}
			query = query.Order("created_at DESC, log_id DESC")
		}
	} else {
		query = query.Order("created_at DESC, log_id DESC")
	}

	err := query.Limit(limit + 1).Find(&logs).Error
	if err != nil {
		return nil, err
	}

	
	if cursor != nil && cursor.Direction != nil && *cursor.Direction == "prev" {
		for i, j := 0, len(logs)-1; i < j; i, j = i+1, j-1 {
			logs[i], logs[j] = logs[j], logs[i]
		}
	}

	return logs, nil
}