package repository

import (
	"Test/internal/feature/permission/entity"
	"context"
	"github.com/google/uuid"
)

type RolePermissionRepository interface {
	AddPermissionToRole(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
	RemovePermissionFromRole(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
	GetPermissionsByRole(ctx context.Context, roleID uuid.UUID) ([]*entity.Permission, error)
}
