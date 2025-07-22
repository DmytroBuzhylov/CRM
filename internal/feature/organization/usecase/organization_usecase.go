package usecase

import (
	"Test/internal/feature/organization/interface_adapters/dto"
	"context"
)

type OrganizationUsecase interface {
	Create(ctx context.Context, req dto.CreateOrganizationRequest) (dto.CreateOrganizationResponse, error)
}
