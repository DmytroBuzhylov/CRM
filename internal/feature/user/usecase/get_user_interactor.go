package usecase

import (
	"Test/internal/feature/user/interface_adapters/dto"
	"Test/internal/feature/user/repository"
	"context"
	"fmt"
)

type getUserInteractor struct {
	userRepo repository.UserRepository
}

func NewGetUserInteractor(userRepo repository.UserRepository) *getUserInteractor {
	return &getUserInteractor{userRepo: userRepo}
}

func (i *getUserInteractor) GetById(ctx context.Context, req dto.GetUserRequest) (dto.GetUserResponse, error) {
	return dto.GetUserResponse{}, fmt.Errorf("")
}
