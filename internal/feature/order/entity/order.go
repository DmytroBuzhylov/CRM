package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Order struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	TableID        uuid.UUID
	UserType       string // 'waiter', 'manager', 'client'
	Status         string // 'pending', 'confirmed', 'completed', 'canceled'
	Notes          string
	TotalPrice     decimal.Decimal
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
