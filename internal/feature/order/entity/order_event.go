package entity

import "github.com/maniartech/signals"

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
