package repository

import (
	"Test/internal/feature/permission/entity"
	"context"
	"github.com/google/uuid"
)

type PermissionRepository interface {
	Create(ctx context.Context, p *entity.Permission) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Permission, error)
	GetByName(ctx context.Context, name string) (*entity.Permission, error)
	GetAll(ctx context.Context) ([]*entity.Permission, error)
	Delete(ctx context.Context, id uuid.UUID) error
	PermissionExists(ctx context.Context, userID uuid.UUID, organizationID uuid.UUID, permissionName string) (bool, error)
}
