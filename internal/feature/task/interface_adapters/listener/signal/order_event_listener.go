package signal

import (
	"Test/internal/feature/order/entity"
	"Test/internal/feature/task/usecase"
	"context"
)

type OrderEventListener struct {
	taskUC usecase.TaskUseCase
}

func NewOrderEventListener(taskUC usecase.TaskUseCase) *OrderEventListener {
	return &OrderEventListener{taskUC: taskUC}
}

func (l *OrderEventListener) Register(ctx context.Context) {
	entity.OrderCreatedEvent.AddListener(func(ctx context.Context, items *entity.OrderWithItems) {
		//"The order has been created. Please contact the guest to clarify and confirm the order."

	})

	entity.OrderCookingEvent.AddListener(func(ctx context.Context, items *entity.OrderWithItems) {

	})

	entity.OrderConfirmedEvent.AddListener(func(ctx context.Context, items *entity.OrderWithItems) {

	})

	entity.OrderCanceledEvent.AddListener(func(ctx context.Context, items *entity.OrderWithItems) {

	})

	entity.OrderPaidEvent.AddListener(func(ctx context.Context, items *entity.OrderWithItems) {

	})

	entity.OrderReadyForPickupEvent.AddListener(func(ctx context.Context, items *entity.OrderWithItems) {

	})

	entity.OrderServedEvent.AddListener(func(ctx context.Context, items *entity.OrderWithItems) {

	})

	entity.OrderUpdatedEvent.AddListener(func(ctx context.Context, items *entity.OrderWithItems) {

	})
}

func (l *OrderEventListener) UnRegister(ctx context.Context) error {
	return nil
}
