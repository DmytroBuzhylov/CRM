package usecase

import (
	"Test/internal/feature/user/interface_adapters/dto"
	"context"
)

type UserUseCase interface {
	Create(ctx context.Context, req dto.CreateUserRequest) (dto.CreateUserResponse, error)
	GetById(ctx context.Context, req dto.GetUserRequest) (dto.GetUserResponse, error)
	GetByUsername(ctx context.Context, req dto.GetUserRequest) (dto.GetUserResponse, error)
	GetByEmail(ctx context.Context, req dto.GetUserRequest) (dto.GetUserResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
	Refresh(ctx context.Context, refreshToken string) (newAccessToken string, newRefreshToken string, err error)
}
