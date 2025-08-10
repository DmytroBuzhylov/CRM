package dto

import (
	"github.com/google/uuid"
	"time"
)

type CreateIngredientRequest struct {
	OrganizationID  uuid.UUID `json:"-"`
	Name            string    `json:"name" binding:"required"`
	Unit            string    `json:"unit" binding:"required"`
	Quantity        uint64    `json:"quantity" binding:"required"`
	MinimumQuantity uint64    `json:"minimum_quantity" binding:"required"`
}

type CreateIngredientResponse struct {
	Status string `json:"status"`
}

type GetIngredientRequest struct {
	ID             uuid.UUID `json:"-"`
	OrganizationID uuid.UUID `json:"-"`
}
type GetIngredientResponse struct {
	ID              uuid.UUID  `json:"id"`
	OrganizationID  uuid.UUID  `json:"-"`
	Name            string     `json:"name"`
	Unit            string     `json:"unit"`
	Quantity        uint64     `json:"quantity"`
	MinimumQuantity uint64     `json:"minimum_quantity"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

type GetAllIngredientsRequest struct {
	OrganizationID uuid.UUID `json:"-"`
}

type DeleteIngredientRequest struct {
	ID             uuid.UUID `json:"-"`
	OrganizationID uuid.UUID `json:"-"`
}

type DeleteManyIngredientsRequest struct {
	IDs            []uuid.UUID `json:"ids"`
	OrganizationID uuid.UUID   `json:"-"`
}

type DeleteIngredientResponse struct {
	Status string `json:"status"`
}

type UpdateIngredientRequest struct {
	ID              uuid.UUID `json:"-"`
	OrganizationID  uuid.UUID `json:"-"`
	Name            string    `json:"name"`
	Unit            string    `json:"unit"`
	Quantity        uint64    `json:"quantity"`
	MinimumQuantity uint64    `json:"minimum_quantity"`
}

type UpdateIngredientResponse struct {
	Status string `json:"status"`
}
