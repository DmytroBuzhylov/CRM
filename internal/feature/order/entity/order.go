package entity

import (
	"github.com/google/uuid"
	"github.com/maniartech/signals"
	"github.com/shopspring/decimal"
	"time"
)

var (
	OrderCreatedEvent        = signals.New[*OrderWithItems]()
	OrderConfirmedEvent      = signals.New[*OrderWithItems]()
	OrderCookingEvent        = signals.New[*OrderWithItems]()
	OrderReadyForPickupEvent = signals.New[*OrderWithItems]()
	OrderServedEvent         = signals.New[*OrderWithItems]()
	OrderPaidEvent           = signals.New[*OrderWithItems]()
	OrderCanceledEvent       = signals.New[*OrderWithItems]()
	OrderUpdatedEvent        = signals.New[*OrderWithItems]()
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
