package dto

import (
	"github.com/google/uuid"
	"time"
)

type MenuItemWithRecipeDTO struct {
	MenuItem    *MenuItemRequest     `json:"menu_item"`
	RecipeItems []*RecipeItemRequest `json:"recipe_items"`
}

type MenuItemRequest struct {
	ID             uuid.UUID `json:"id,omitempty"`
	OrganizationID uuid.UUID `json:"organization_id,omitempty"`
	Name           string    `json:"name" binding:"required,min=8"`
	Description    string    `json:"description" binding:"required"`
	Price          string    `json:"price" binding:"required"`
	Category       string    `json:"category" binding:"required"`
	IsAvailable    bool      `json:"is_available,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
}

type RecipeItemRequest struct {
	ID             uuid.UUID `json:"id,omitempty"`
	MenuItemID     uuid.UUID `json:"menu_item_id,omitempty"`
	IngredientID   uuid.UUID `json:"ingredient_id"`
	QuantityNeeded uint64    `json:"quantity_needed" binding:"required"`
	UnitOfMeasure  string    `json:"unit_of_measure" binding:"required"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
}

type CreateMenuItemRequest struct {
	OrganizationID uuid.UUID `json:"-"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Price          string    `json:"price"`
	Category       string    `json:"category"`
}

//type CreateMenuItemResponse struct {
//	Status string `json:"status"`
//}

type GetMenuItemRequest struct {
	ID uuid.UUID `json:"id"`
}

type GetMenuItemResponse struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	Category    string    `json:"category"`
	IsAvailable bool      `json:"is_available"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type GetAllMenuItemsRequest struct {
	OrganizationID uuid.UUID `json:"-"`
}

//type UpdateMenuItemRequest struct {
//	ID          uuid.UUID `json:"id"`
//	Name        string    `json:"name"`
//	Price       float64   `json:"price"`
//	Category    string    `json:"category"`
//	Description string    `json:"description"`
//}

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
