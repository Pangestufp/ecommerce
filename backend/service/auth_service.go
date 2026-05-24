package service

import (
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"backend/helper"
	"backend/repository"
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthService interface {
	Register(req *dto.RegisterRequest, userType string) error
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, *string, error)
	Refresh(ctx context.Context, userRefreshToken string, userID string) (*dto.LoginResponse, *string, error)
}

type authService struct {
	repositoryA repository.AuthRepository
	redis       *redis.Client
}

func NewAuthService(repositoryA repository.AuthRepository, redis *redis.Client) *authService {
	return &authService{
		repositoryA: repositoryA,
		redis:       redis,
	}
}

func (s *authService) Register(req *dto.RegisterRequest, userType string) error {

	if !strings.Contains(req.Email, "@") || !strings.HasSuffix(req.Email, ".com") {
		return &errorhandler.BadRequestError{Message: "invalid email format"}
	}

	if emailExist := s.repositoryA.EmailExist(helper.LowerAndTrim(req.Email)); emailExist {
		return &errorhandler.BadRequestError{Message: "email already registered"}
	}

	if len(req.Password) < 10 {
		return &errorhandler.BadRequestError{Message: "password to short"}
	}

	if req.Password != req.PasswordConfirmation {
		return &errorhandler.BadRequestError{Message: "password not match"}
	}

	passwordHash, err := helper.HashPassword(req.Password)

	if err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	user := entity.User{
		UserID:     uuid.New().String(),
		Name:       helper.UpperAndTrim(req.Name),
		Email:      helper.LowerAndTrim(req.Email),
		Password:   passwordHash,
		CreatedAt:  helper.TimeNowWIB(),
		UpdatedAt:  helper.TimeNowWIB(),
		Status:     1,
		Role:       userType,
		VerifiedAt: nil,
	}

	if err := s.repositoryA.Register(&user); err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	return nil
}

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, *string, error) {

	user, err := s.repositoryA.GetUserByEmail(helper.LowerAndTrim(req.Email))

	if err != nil {
		return nil, nil, &errorhandler.NotFoundError{Message: "wrong email or password"}
	}

	if err := helper.VerifyPassword(user.Password, req.Password); err != nil {
		return nil, nil, &errorhandler.NotFoundError{Message: "wrong email or password"}
	}

	token, err := helper.GenerateToken(user)

	if err != nil {
		return nil, nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	refreshToken, err := helper.GenerateRefreshToken()
	if err != nil {
		return nil, nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	now := time.Now()
	expiredAt := now.Add(7 * 24 * time.Hour)

	refreshTokenHash, err := helper.HashRefreshToken(refreshToken)

	if err != nil {
		return nil, nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	newRefreshToken := entity.RefreshToken{
		ID:               uuid.New().String(),
		UserID:           user.UserID,
		RefreshTokenHash: refreshTokenHash,
		ExpiredAt:        expiredAt,
		CreatedAt:        helper.TimeNowWIB(),
		UpdatedAt:        helper.TimeNowWIB(),
	}

	if err := s.repositoryA.StoreRefreshToken(ctx, &newRefreshToken); err != nil {
		return nil, nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	res := dto.LoginResponse{
		Token: token,
	}

	return &res, &refreshToken, nil
}

func (s *authService) Refresh(ctx context.Context, userRefreshToken string, userID string) (*dto.LoginResponse, *string, error) {
	user, err := s.repositoryA.GetUserByID(userID) //cek user
	if err != nil {
		return nil, nil, &errorhandler.UnauthorizedError{Message: "You should login"}
	}

	now := time.Now()

	oldRefreshToken, err := s.repositoryA.GetRefreshToken(ctx, userID, now) //ambil token lama yg udh dihash
	if err != nil {
		return nil, nil, &errorhandler.UnauthorizedError{Message: "You should login"}
	}

	if err := helper.VerifyRefreshToken(oldRefreshToken.RefreshTokenHash, userRefreshToken); err != nil { //cek token lama sama gak ama yg dihit
		return nil, nil, &errorhandler.UnauthorizedError{Message: "Your Token is Invalid"}
	}

	token, err := helper.GenerateToken(user) //buat jwt baru

	if err != nil {
		return nil, nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	refreshToken, err := helper.GenerateRefreshToken() //buat refreshtoken baru
	if err != nil {
		return nil, nil, &errorhandler.InternalServerError{Message: err.Error()}
	}
	hashRefreshToken, err := helper.HashRefreshToken(refreshToken) //hash refrresh tokennya

	if err != nil {
		return nil, nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	newRefreshToken := entity.RefreshToken{
		ID:               uuid.New().String(),
		UserID:           user.UserID,
		RefreshTokenHash: hashRefreshToken,
		ExpiredAt:        oldRefreshToken.ExpiredAt,
		CreatedAt:        oldRefreshToken.CreatedAt,
		UpdatedAt:        helper.TimeNowWIB(),
	}

	if err := s.repositoryA.UpdateRefreshToken(ctx, &newRefreshToken); err != nil { //buat token baru
		return nil, nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	res := dto.LoginResponse{ //kasih respons
		Token: token,
	}

	return &res, &refreshToken, nil
}
