package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type OrderItem struct {
	ID                 uuid.UUID
	OrderID            uuid.UUID
	MenuItemID         uuid.UUID
	Quantity           uint64
	PriceAtTimeOfOrder decimal.Decimal
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
