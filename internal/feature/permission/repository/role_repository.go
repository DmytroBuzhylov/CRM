package repository

import (
	"Test/internal/feature/permission/entity"
	"context"
	"github.com/google/uuid"
)

type RoleRepository interface {
	Create(ctx context.Context, role *entity.Role) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Role, error)
	GetMany(ctx context.Context, organizationID uuid.UUID) ([]*entity.Role, error)
	Update(ctx context.Context, role *entity.Role) error
	Delete(ctx context.Context, id uuid.UUID) error
}
