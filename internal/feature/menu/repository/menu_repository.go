package repository

import (
	"Test/internal/feature/menu/entity"
	"Test/internal/feature/menu/repository/postgres"
	"context"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type MenuRepository interface {
	GetMenuItem(ctx context.Context, menuItemID uuid.UUID) (item entity.MenuItem, err error)
	GetAllMenuItems(ctx context.Context, organizationID uuid.UUID) (items []entity.MenuItem, err error)
}

type MenuTransactionalRepository interface {
	WithTx(ctx context.Context) (*postgres.PostgresMenuTransactionalRepository, error)

	CreateMenuItem(ctx context.Context, menuItem entity.MenuItem) error
	UpdateMenuItem(ctx context.Context, menuItem entity.MenuItem) error
	DeleteMenuItem(ctx context.Context, menuItemID uuid.UUID, organizationID uuid.UUID) error
	AddRecipeItem(ctx context.Context, recipeItem entity.RecipeItem) error
	AddRecipeItems(ctx context.Context, recipeItem []entity.RecipeItem) error
	UpdateRecipeItem(ctx context.Context, recipeItem entity.RecipeItem) error
	DeleteRecipeItem(ctx context.Context, itemID uuid.UUID, organizationID uuid.UUID) error

	DecreaseInventory(ctx context.Context, quantity decimal.NullDecimal, ingredientID uuid.UUID, organizationID uuid.UUID) error
}
