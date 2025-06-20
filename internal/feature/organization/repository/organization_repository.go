package repository

import (
	"Test/internal/feature/organization/entity"
	"context"
)

type OrganizationRepository interface {
	Create(ctx context.Context, organization entity.Organization) error
	UpdateName(ctx context.Context, organizationID uint64, name string) error
	UpdateDescription(ctx context.Context, organizationID uint64, description string) error
	UpdateUsers(ctx context.Context, organizationID uint64, userID uint64) error
}
