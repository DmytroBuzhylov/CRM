package usecase

import (
	"Test/config"
	"Test/internal/feature/user/entity"
	"Test/internal/feature/user/interface_adapters/dto"
	"Test/internal/feature/user/repository"
	"Test/internal/feature/user/repository/postgres"
	"Test/internal/pkg/jwt"
	"context"
	"errors"
	"time"
)

type createUserInteractor struct {
	userRepo  repository.UserRepository
	jwtConfig config.JWTConfig
}

func NewCreateUserInteractor(repo repository.UserRepository, jwtConfig config.JWTConfig) *createUserInteractor {
	return &createUserInteractor{userRepo: repo, jwtConfig: jwtConfig}
}

func (i *createUserInteractor) Create(ctx context.Context, req dto.CreateUserRequest) (dto.CreateUserResponse, error) {
	hashedPassword, err := postgres.HashPassword(req.Password)
	if err != nil {
		return dto.CreateUserResponse{}, errors.New("error hashing password")
	}
	user := entity.NewUser(
		req.FullName,
		req.Username,
		hashedPassword,
		req.Email,
		req.Phone,
		"platform_user",
	)

	err = i.userRepo.Create(ctx, user)
	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	newAccessToken, newRefreshToken, err := jwt.GenerateTokens(
		user.ID,
		user.Role,
		nil,
		i.jwtConfig.JWTAccessSecret,
		i.jwtConfig.JWTRefreshSecret,
		i.jwtConfig.JWTAccessLifetime,
		i.jwtConfig.JWTRefreshLifetime,
	)
	if err != nil {
		return dto.CreateUserResponse{}, errors.New("error generate jwt tokens")
	}

	err = i.userRepo.SaveRefreshToken(ctx, user.ID, newRefreshToken, time.Now().Add(i.jwtConfig.JWTRefreshLifetime))
	if err != nil {
		_ = i.userRepo.DeleteUser(ctx, user.ID)
		return dto.CreateUserResponse{}, err
	}

	return dto.CreateUserResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil

}
