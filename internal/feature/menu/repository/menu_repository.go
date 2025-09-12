package repository

import (
	"Test/internal/feature/menu/entity"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type MenuRepository interface {
	GetMenuItem(ctx context.Context, menuItemID uuid.UUID) (item entity.MenuItem, err error)
	GetAllMenuItems(ctx context.Context, organizationID uuid.UUID) (items []entity.MenuItem, err error)
	StartTx(ctx context.Context, opts pgx.TxOptions) (MenuTransactionalRepository, error)
}

type MenuTransactionalRepository interface {
	Rollback(ctx context.Context) error
	Commit(ctx context.Context) error

	CreateMenuItem(ctx context.Context, menuItem entity.MenuItem) error
	UpdateMenuItem(ctx context.Context, menuItem entity.MenuItem) error
	DeleteMenuItem(ctx context.Context, menuItemID uuid.UUID, organizationID uuid.UUID) error
	AddRecipeItem(ctx context.Context, recipeItem entity.RecipeItem) error
	AddRecipeItems(ctx context.Context, recipeItem []entity.RecipeItem) error
	UpdateRecipeItem(ctx context.Context, recipeItem entity.RecipeItem) error
	DeleteRecipeItem(ctx context.Context, itemID uuid.UUID, organizationID uuid.UUID) error

	DecreaseInventory(ctx context.Context, menuItemID uuid.UUID, quantity uint, organizationID uuid.UUID) error
}
