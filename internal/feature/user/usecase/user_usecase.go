package usecase

import (
	"Test/internal/feature/user/interface_adapters/dto"
	"context"
)

type CreateUser interface {
	Create(ctx context.Context, req dto.CreateUserRequest) (dto.CreateUserResponse, error)
}

type GetUser interface {
	Get(ctx context.Context, req dto.GetUserRequest) (dto.GetUserResponse, error)
}

type UpdateUser interface {
	Update(ctx context.Context, req dto.UpdateUserRequest) (dto.UpdateUserResponse, error)
}

type DeleteUser interface {
	Delete(ctx context.Context, req dto.DeleteUserRequest) (dto.DeleteUserResponse, error)
}
