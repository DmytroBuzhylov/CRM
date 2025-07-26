package usecase

import (
	"Test/internal/feature/organization/entity"
	"Test/internal/feature/organization/interface_adapters/dto"
	"Test/internal/feature/organization/repository"
	"context"
)

type organizationUseCaseInteractor struct {
	repo repository.OrganizationRepository
}

func NewOrganizationUseCaseInteractor(repo repository.OrganizationRepository) *organizationUseCaseInteractor {
	return &organizationUseCaseInteractor{repo: repo}
}

func (uc *organizationUseCaseInteractor) Create(ctx context.Context, req dto.CreateOrganizationRequest) (dto.CreateOrganizationResponse, error) {
	newOrganization := entity.NewOrganization(req.Name, req.Description, req.OwnerUserID)

	err := uc.repo.Create(ctx, newOrganization)
	if err != nil {
		return dto.CreateOrganizationResponse{}, err
	}
	return dto.CreateOrganizationResponse{Status: "ok"}, nil
}
