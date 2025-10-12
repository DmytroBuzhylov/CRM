package repository

import (
	"Test/internal/feature/order/entity"
	"context"
	"github.com/google/uuid"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *entity.Order, items []*entity.OrderItem) error
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, organizationID uuid.UUID, status string) error
	UpdateOrderNotes(ctx context.Context, id uuid.UUID, organizationID uuid.UUID, notes string) error
	GetOrder(ctx context.Context, id uuid.UUID) (*entity.Order, error)
	GetAllOrders(ctx context.Context, organizationID uuid.UUID) ([]*entity.Order, error)

	GetOrderWithItems(ctx context.Context, id uuid.UUID) (*entity.OrderWithItems, error)
	GetAllOrdersWithItems(ctx context.Context, organizationID uuid.UUID) ([]*entity.OrderWithItems, error)

	GetOrderItem(ctx context.Context, id uuid.UUID) (*entity.OrderItem, error)
	GetOrderItems(ctx context.Context, id []uuid.UUID) ([]*entity.OrderItem, error)

	CreateOrderItems(ctx context.Context, organizationID uuid.UUID, orderID uuid.UUID, items []*entity.OrderItem) error
}
