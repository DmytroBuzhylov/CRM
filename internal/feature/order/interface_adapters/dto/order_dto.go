package dto

import (
	"Test/internal/feature/order/entity"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type OrderWithOrderItemDTO struct {
	Order      OrderResponse       `json:"order"`
	OrderItems []OrderItemResponse `json:"order_items"`
}

type OrderResponse struct {
	ID             uuid.UUID       `json:"id"`
	OrganizationID uuid.UUID       `json:"organization_id"`
	UserID         uuid.UUID       `json:"user_id"`
	TableID        uuid.UUID       `json:"table_id"`
	UserType       string          `json:"user_type"`
	Status         string          `json:"status"`
	Notes          string          `json:"notes"`
	TotalPrice     decimal.Decimal `json:"total_price"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type OrderItemResponse struct {
	ID                 uuid.UUID       `json:"id"`
	OrderID            uuid.UUID       `json:"order_id"`
	MenuItemID         uuid.UUID       `json:"menu_item_id"`
	Quantity           uint64          `json:"quantity"`
	PriceAtTimeOfOrder decimal.Decimal `json:"price_at_time_of_order"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

func OrderWithItemsToDTO(o *entity.OrderWithItems) OrderWithOrderItemDTO {
	items := make([]OrderItemResponse, 0, len(o.Items))
	for _, item := range o.Items {
		items = append(items, OrderItemResponse{
			ID:                 item.ID,
			OrderID:            item.OrderID,
			MenuItemID:         item.MenuItemID,
			Quantity:           item.Quantity,
			PriceAtTimeOfOrder: item.PriceAtTimeOfOrder,
			CreatedAt:          item.CreatedAt,
			UpdatedAt:          item.UpdatedAt,
		})
	}

	return OrderWithOrderItemDTO{
		Order: OrderResponse{
			ID:             o.Order.ID,
			OrganizationID: o.Order.OrganizationID,
			UserID:         o.Order.UserID,
			TableID:        o.Order.TableID,
			UserType:       o.Order.UserType,
			Status:         o.Order.Status,
			Notes:          o.Order.Notes,
			TotalPrice:     o.Order.TotalPrice,
			CreatedAt:      o.Order.CreatedAt,
			UpdatedAt:      o.Order.UpdatedAt,
		},
		OrderItems: items,
	}
}

type CreateOrderRequest struct {
	Order struct {
		OrganizationID uuid.UUID `json:"-"`
		UserID         uuid.UUID `json:"userID"`
		TableID        uuid.UUID `json:"tableID" binding:"required"`
		UserType       string    `json:"user_type"`
		Notes          string    `json:"notes"`
	} `json:"order"`
	OrderItems []struct {
		MenuItemID uuid.UUID `json:"menu_item_id" binding:"required"`
		Quantity   uint64    `json:"quantity" binding:"required"`
	} `json:"order_items"`
}

type GetOrderRequest struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

type GetOrganizationOrdersRequest struct {
	OrganizationID uuid.UUID `json:"-"`
}

type UpdateOrderStatusRequest struct {
	ID             uuid.UUID `json:"id" binding:"required"`
	OrganizationID uuid.UUID `json:"-"`
	Status         string    `json:"status" binding:"required"`
}

type UpdateOrderNotesRequest struct {
	ID             uuid.UUID `json:"id" binding:"required"`
	OrganizationID uuid.UUID `json:"-"`
	Notes          string    `json:"notes" binding:"required"`
}

type AddOrderItemsToOrderRequest struct {
	OrganizationID uuid.UUID `json:"-"`
	OrderID        uuid.UUID `json:"order_id" binding:"required"`
	OrderItems     []struct {
		MenuItemID uuid.UUID `json:"menu_item_id" binding:"required"`
		Quantity   uint64    `json:"quantity" binding:"required"`
	}
}
