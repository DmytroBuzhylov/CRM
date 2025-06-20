package usecase

import (
	"Test/internal/feature/user/interface_adapters/dto"
	"Test/internal/feature/user/repository"
	"context"
)

type createUserInteractor struct {
	userRepo repository.UserRepository
}

func NewCreateUserInteractor(repo repository.UserRepository) *createUserInteractor {
	return &createUserInteractor{userRepo: repo}
}

func (i *createUserInteractor) Create(ctx context.Context, req dto.CreateUserRequest) (dto.CreateUserResponse, error) {
	
}
