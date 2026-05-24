package repository

import (
	"backend/entity"
	"context"
	"time"

	"gorm.io/gorm"
)

type AuthRepository interface {
	EmailExist(email string) bool
	Register(req *entity.User) error
	GetUserByEmail(email string) (*entity.User, error)
	GetRefreshToken(ctx context.Context, userID string, now time.Time) (*entity.RefreshToken, error)
	StoreRefreshToken(ctx context.Context, token *entity.RefreshToken) error
	UpdateRefreshToken(ctx context.Context, token *entity.RefreshToken) error
	GetUserByID(userID string) (*entity.User, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *authRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) EmailExist(email string) bool {
	var user entity.User

	err := r.db.First(&user, "email = ? and status = 1", email).Error

	return err == nil
}

func (r *authRepository) Register(user *entity.User) error {
	err := r.db.Create(&user).Error

	return err
}

func (r *authRepository) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User

	err := r.db.First(&user, "email = ?", email).Error

	return &user, err
}

func (r *authRepository) GetRefreshToken(ctx context.Context, userID string, now time.Time) (*entity.RefreshToken, error) {
	var refreshToken entity.RefreshToken

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND expired_at >= ?", userID, now).
		First(&refreshToken).Error

	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *authRepository) StoreRefreshToken(ctx context.Context, token *entity.RefreshToken) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.
			Where("user_id = ?", token.UserID).
			Delete(&entity.RefreshToken{}).Error; err != nil {
			return err
		}

		if err := tx.Create(token).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *authRepository) GetUserByID(userID string) (*entity.User, error) {
	var user entity.User

	err := r.db.First(&user, "user_id = ?", userID).Error

	return &user, err
}

func (r *authRepository) UpdateRefreshToken(ctx context.Context, token *entity.RefreshToken) error {
	return r.db.WithContext(ctx).
		Model(&entity.RefreshToken{}).
		Where("user_id = ?", token.UserID).
		Updates(map[string]interface{}{
			"refresh_token_hash": token.RefreshTokenHash,
			"expired_at":         token.ExpiredAt,
			"updated_at":         token.UpdatedAt,
		}).Error
}
