package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type MenuItem struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Name           string
	Description    string
	Price          decimal.NullDecimal `db:"price"`
	Category       string
	IsAvailable    bool
	UpdatedAt      time.Time
	CreatedAt      time.Time
}
