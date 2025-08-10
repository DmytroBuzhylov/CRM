package repository

import (
	"Test/internal/feature/menu/entity"
	"context"
	"github.com/google/uuid"
)

type MenuRepository interface {
	CreateMenuItem(ctx context.Context, item entity.MenuItem) error
	GetMenuItem(ctx context.Context, id uuid.UUID) (item entity.MenuItem, err error)
	GetAllMenuItems(ctx context.Context, organizationID uuid.UUID) (items []entity.MenuItem, err error)
	UpdateMenuItem(ctx context.Context, item entity.MenuItem) error
	DeleteMenuItem(ctx context.Context, itemID uuid.UUID, organizationID uuid.UUID) error

	AddRecipeItem(ctx context.Context)
	UpdateRecipeItem(ctx context.Context)
	DeleteRecipeItem(ctx context.Context)

	DecreaseInventory(ctx context.Context)
}
