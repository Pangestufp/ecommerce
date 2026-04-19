package service

import (
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"backend/helper"
	"backend/repository"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthService interface {
	Register(req *dto.RegisterRequest, userType string) error
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
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
		Address:    req.Address,
		Phone:      req.Phone,
		PostalCode: req.PostalCode,
		VerifiedAt: nil,
	}

	if err := s.repositoryA.Register(&user); err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	return nil
}

func (s *authService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {

	user, err := s.repositoryA.GetUserByEmail(helper.LowerAndTrim(req.Email))

	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "wrong email or password"}
	}

	if err := helper.VerifyPassword(user.Password, req.Password); err != nil {
		return nil, &errorhandler.NotFoundError{Message: "wrong email or password"}
	}

	token, err := helper.GenerateToken(user)

	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	res := dto.LoginResponse{
		ID:    user.UserID,
		Name:  user.Name,
		Token: token,
	}

	return &res, err
}
