package usecase

import (
	"Test/config"
	"Test/internal/feature/user/entity"
	"Test/internal/feature/user/interface_adapters/dto"
	"Test/internal/feature/user/repository"
	"Test/internal/feature/user/repository/postgres"
	"Test/internal/pkg/jwt"
	"context"
	"fmt"
)

type getUserInteractor struct {
	userRepo  repository.UserRepository
	jwtConfig config.JWTConfig
}

func NewGetUserInteractor(userRepo repository.UserRepository) *getUserInteractor {
	return &getUserInteractor{userRepo: userRepo}
}

func (i *getUserInteractor) GetById(ctx context.Context, req dto.GetUserRequest) (dto.GetUserResponse, error) {
	return dto.GetUserResponse{}, fmt.Errorf("")
}

func (i *getUserInteractor) GetByUsername(ctx context.Context, req dto.GetUserRequest) (dto.GetUserResponse, error) {

}

func (i *getUserInteractor) GetByEmail(ctx context.Context, req dto.GetUserRequest) (dto.GetUserResponse, error) {

}

func (i *getUserInteractor) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	var (
		user entity.User
		err  error
	)
	switch req.Method {
	case "email":
		user, err = i.userRepo.FindByEmail(ctx, req.Email)

	case "phone":
		user, err = i.userRepo.FindByPhone(ctx, req.Phone)

	case "username":
		user, err = i.userRepo.FindByUsername(ctx, req.Username)

	default:
		return dto.LoginResponse{Status: "incorrect method"}, fmt.Errorf("incorrect method")
	}

	if err != nil {
		return dto.LoginResponse{Status: "db error"}, err
	}
	if !postgres.CheckPasswordHash(req.Password, user.HashedPassword) {
		return dto.LoginResponse{Status: "incorrect password"}, err
	}

	accesToken, refreshToken, err := jwt.GenerateTokens(user.ID, user.Role, i.jwtConfig.JWTAccessSecret, i.jwtConfig.JWTRefreshSecret, i.jwtConfig.JWTAccessLifetime, i.jwtConfig.JWTRefreshLifetime)

	if err != nil {
		return dto.LoginResponse{Status: "JWT generation error"}, err
	}

}

func (i *getUserInteractor) Refresh(ctx context.Context, refreshToken string) error {
	claims, err := jwt.VerifyToken(refreshToken, i.jwtConfig.JWTRefreshSecret)
	if err != nil {
		return err
	}
	
}
