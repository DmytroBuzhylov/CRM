package usecase

import (
	"Test/internal/feature/user/interface_adapters/dto"
	"context"
)

type CreateUser interface {
	Create(ctx context.Context, req dto.CreateUserRequest) (dto.CreateUserResponse, error)
}

type GetUser interface {
	GetById(ctx context.Context, req dto.GetUserRequest) (dto.GetUserResponse, error)
	GetByUsername(ctx context.Context, req dto.GetUserRequest) (dto.GetUserResponse, error)
	GetByEmail(ctx context.Context, req dto.GetUserRequest) (dto.GetUserResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
}

type UpdateUser interface {
	Update(ctx context.Context, req dto.UpdateUserRequest) (dto.UpdateUserResponse, error)
}

type DeleteUser interface {
	Delete(ctx context.Context, req dto.DeleteUserRequest) (dto.DeleteUserResponse, error)
}
