package usecase

import (
	"Test/internal/feature/permission/entity"
	"context"
	"github.com/google/uuid"
)

type PermissionUseCase interface {
	RequestUserPermissions(ctx context.Context, userID uuid.UUID) ([]*entity.Permission, error)
	HasPermission(ctx context.Context, userID uuid.UUID, organizationID uuid.UUID, permissionName string) bool
	AddPermissions(ctx context.Context, roleID uuid.UUID, permissionsID []uuid.UUID) error
	DeletePermissions(ctx context.Context, roleID uuid.UUID, permissionsID []uuid.UUID) error
}
