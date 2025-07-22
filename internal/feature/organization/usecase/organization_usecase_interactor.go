package usecase

import (
	"Test/internal/feature/organization/entity"
	"Test/internal/feature/organization/interface_adapters/dto"
	"Test/internal/feature/organization/repository"
	"context"
)

type organizationUsecaseInteractor struct {
	repo repository.OrganizationRepository
}

func NewOrganizationUsecaseInteractor(repo repository.OrganizationRepository) *organizationUsecaseInteractor {
	return &organizationUsecaseInteractor{repo: repo}
}

func (uc *organizationUsecaseInteractor) Create(ctx context.Context, req dto.CreateOrganizationRequest) (dto.CreateOrganizationResponse, error) {
	newOrganization := entity.NewOrganization(req.Name, req.Description, req.OwnerUserID)

	err := uc.repo.Create(ctx, newOrganization)
	if err != nil {
		return dto.CreateOrganizationResponse{}, err
	}
	return dto.CreateOrganizationResponse{Status: "ok"}, nil
}
