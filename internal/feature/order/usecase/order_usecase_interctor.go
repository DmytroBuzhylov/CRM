package usecase

import (
	"Test/internal/feature/order/entity"
	"Test/internal/feature/order/interface_adapters/dto"
	"Test/internal/feature/order/repository"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type OrderUseCaseInteractor struct {
	repo repository.OrderRepository
	//evBus events.EventBus
}

func NewOrderUseCaseInteractor(repo repository.OrderRepository) *OrderUseCaseInteractor {
	return &OrderUseCaseInteractor{repo: repo}
}

func (uc *OrderUseCaseInteractor) CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) error {
	newOrder := &entity.Order{
		ID:             uuid.New(),
		OrganizationID: req.Order.OrganizationID,
		UserID:         req.Order.UserID,
		TableID:        req.Order.TableID,
		UserType:       req.Order.UserType,
		Notes:          req.Order.Notes,
	}

	orderItems := make([]*entity.OrderItem, 0, len(req.OrderItems))
	for _, i := range req.OrderItems {
		orderItems = append(orderItems, &entity.OrderItem{
			MenuItemID: i.MenuItemID,
			Quantity:   i.Quantity,
		})
	}

	err := uc.repo.CreateOrder(ctx, newOrder, orderItems)
	if err != nil {
		log.Warn().Err(err).
			Str("usecase", "Order").
			Str("method", "CreateOrder").
			Msg("failed to create order")
		return errors.New("failed to create order")
	}

	entity.OrderCreatedEvent.Emit(ctx, &entity.OrderWithItems{
		Order: newOrder,
		Items: orderItems,
	})

	return nil
}

func (uc *OrderUseCaseInteractor) GetOrder(ctx context.Context, req *dto.GetOrderRequest) (*dto.OrderWithOrderItemDTO, error) {

	orderWithItems, err := uc.repo.GetOrderWithItems(ctx, req.ID)
	if err != nil {
		if errors.As(err, pgx.ErrNoRows) {
			return nil, errors.New("this order does not exist")
		}
		log.Warn().Err(err).
			Str("usecase", "Order").
			Str("method", "GetOrder").
			Msg("failed to get order")
		return nil, errors.New("failed to get order")
	}

	orderItems := make([]*entity.OrderItem, 0, len(orderWithItems.Items))
	for _, i := range orderWithItems.Items {
		orderItems = append(orderItems, i)
	}

	resp := dto.OrderWithItemsToDTO(orderWithItems)

	return &resp, nil
}

func (uc *OrderUseCaseInteractor) GetOrganizationOrders(ctx context.Context, req *dto.GetOrganizationOrdersRequest) ([]*dto.OrderWithOrderItemDTO, error) {

	ordersWithItems, err := uc.repo.GetAllOrdersWithItems(ctx, req.OrganizationID)
	if err != nil {
		if errors.As(err, pgx.ErrNoRows) {
			return nil, errors.New("this orders does not exist")
		}
		log.Warn().Err(err).
			Str("usecase", "Order").
			Str("method", "GetOrganizationOrders").
			Msg("failed to get organization orders")
		return nil, errors.New("failed to get organization orders")
	}

	respDTO := make([]*dto.OrderWithOrderItemDTO, 0, len(ordersWithItems))
	for _, item := range ordersWithItems {
		resp := dto.OrderWithItemsToDTO(item)
		respDTO = append(respDTO, &resp)
	}

	return respDTO, nil
}

func (uc *OrderUseCaseInteractor) UpdateOrderStatus(ctx context.Context, req *dto.UpdateOrderStatusRequest) error {

	err := uc.repo.UpdateOrderStatus(ctx, req.ID, req.OrganizationID, req.Status)
	if err != nil {
		log.Warn().Err(err).
			Str("usecase", "Order").
			Str("method", "UpdateOrderStatus").
			Msg("failed to update order status")
		return errors.New("failed to update order status")
	}

	return nil
}

func (uc *OrderUseCaseInteractor) UpdateOrderNotes(ctx context.Context, req *dto.UpdateOrderNotesRequest) error {

	err := uc.repo.UpdateOrderNotes(ctx, req.ID, req.OrganizationID, req.Notes)
	if err != nil {
		log.Warn().Err(err).
			Str("usecase", "Order").
			Str("method", "UpdateOrderNotes").
			Msg("failed to update order notes")
		return errors.New("failed to update order notes")
	}

	return nil
}

func (uc *OrderUseCaseInteractor) AddOrderItemsToOrder(ctx context.Context, req *dto.AddOrderItemsToOrderRequest) error {

	orderItems := make([]*entity.OrderItem, 0, len(req.OrderItems))
	for _, i := range req.OrderItems {
		orderItems = append(orderItems, &entity.OrderItem{
			MenuItemID: i.MenuItemID,
			Quantity:   i.Quantity,
		})
	}

	err := uc.repo.CreateOrderItems(ctx, req.OrganizationID, req.OrderID, orderItems)
	if err != nil {
		log.Warn().Err(err).
			Str("usecase", "Order").
			Str("method", "AddOrderItemsToOrder").
			Msg("failed to add order items to order")
		return errors.New("failed to add order items to order")
	}

	return nil
}
