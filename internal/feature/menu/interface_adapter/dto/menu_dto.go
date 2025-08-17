package dto

import (
	"Test/internal/feature/menu/entity"
	"github.com/google/uuid"
	"time"
)

type MenuItemWithRecipeDTO struct {
	MenuItem    entity.MenuItem
	RecipeItems []entity.RecipeItem
}

type CreateMenuItemRequest struct {
	OrganizationID uuid.UUID `json:"-"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Price          float64   `json:"price"`
	Category       string    `json:"category"`
}

type CreateMenuItemResponse struct {
	Status string `json:"status"`
}

type GetMenuItemRequest struct {
	ID uuid.UUID `json:"id"`
}

type GetMenuItemResponse struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"`
	IsAvailable bool      `json:"is_available"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type GetAllMenuItemsRequest struct {
	OrganizationID uuid.UUID `json:"-"`
}

type UpdateMenuItemRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
}

type UpdateMenuItemResponse struct {
	Status string `json:"status"`
}

type DeleteMenuItemRequest struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"-"`
}

type DeleteMenuItemResponse struct {
	Status string `json:"status"`
}

type AddRecipeItemRequest struct {
	MenuItemID     uuid.UUID `json:"menu_item_id"`
	IngredientID   uuid.UUID `json:"ingredient_id"`
	QuantityNeeded uint64    `json:"quantity_needed"`
	UnitOfMeasure  string    `json:"unit_of_measure"`

	OrganizationID uuid.UUID `json:"-"`
}

type AddRecipeItemResponse struct {
	Status string `json:"status"`
}

type UpdateRecipeItemRequest struct {
	ID             uuid.UUID `json:"id"`
	UnitOfMeasure  string    `json:"unit_of_measure"`
	QuantityNeeded uint64    `json:"quantity_needed"`

	OrganizationID uuid.UUID `json:"-"`
}

type UpdateRecipeItemResponse struct {
	Status string `json:"status"`
}

type DeleteRecipeItemRequest struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"-"`
}

type DeleteRecipeItemResponse struct {
	Status string `json:"status"`
}
