package entity

import (
	"github.com/google/uuid"
	"time"
)

type RecipeItem struct {
	ID             uuid.UUID
	MenuItemID     uuid.UUID
	IngredientID   uuid.UUID
	QuantityNeeded uint64
	UnitOfMeasure  string
	UpdatedAt      time.Time
	CreatedAt      time.Time
}
