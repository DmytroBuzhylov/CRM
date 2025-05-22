package usecase

import (
	"Test/internal/feature/user/interface_adapters/dto"
	"Test/internal/feature/user/repository"
	"context"
)

type getUserInteractor struct {
	userRepo repository.UserRepository
}

func NewGetUserInteractor(userRepo repository.UserRepository) *getUserInteractor {
	return &getUserInteractor{userRepo: userRepo}
}

func (i *getUserInteractor) Get(ctx context.Context, req dto.GetUserRequest) (dto.GetUserResponse, error) {
	i.userRepo.
}
