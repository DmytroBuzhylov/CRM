package usecase

import (
	"Test/internal/feature/order/interface_adapters/dto"
	"context"
)

type OrderUseCase interface {
	CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) error
	GetOrder(ctx context.Context, req *dto.GetOrderRequest) (*dto.OrderWithOrderItemDTO, error)
	GetOrganizationOrders(ctx context.Context, req *dto.GetOrganizationOrdersRequest) ([]*dto.OrderWithOrderItemDTO, error)
	UpdateOrderStatus(ctx context.Context, req *dto.UpdateOrderStatusRequest) error
	UpdateOrderNotes(ctx context.Context, req *dto.UpdateOrderNotesRequest) error
	AddOrderItemsToOrder(ctx context.Context, req *dto.AddOrderItemsToOrderRequest) error
}
