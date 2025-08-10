package entity

import (
	"github.com/google/uuid"
	"time"
)

type Ingredient struct {
	ID              uuid.UUID
	OrganizationID  uuid.UUID
	Name            string
	Unit            string
	Quantity        uint64
	MinimumQuantity uint64
	UpdatedAt       *time.Time
	CreatedAt       *time.Time
}
