package repository

import (
	"Test/internal/feature/inventory/entity"
	"context"
	"github.com/google/uuid"
)

type InventoryRepository interface {
	Get(ctx context.Context, id uuid.UUID, organizationID uuid.UUID) (entity.Ingredient, error)
	GetByName(ctx context.Context, name string) (entity.Ingredient, error)
	GetAll(ctx context.Context, organizationID uuid.UUID) ([]entity.Ingredient, error)
	Create(ctx context.Context, ingredient entity.Ingredient) error
	CreateMany(ctx context.Context, ingredients []entity.Ingredient) error
	Update(ctx context.Context, ingredient entity.Ingredient) error
	Delete(ctx context.Context, id uuid.UUID, organizationID uuid.UUID) error
	DeleteMany(ctx context.Context, ids []uuid.UUID, organizationID uuid.UUID) error
}
