package dto

import (
	"Test/internal/feature/order/entity"
	"github.com/google/uuid"
)

type OrderWithOrderItemDTO struct {
	Order      entity.Order
	OrderItems []entity.OrderItem
}

type CreateOrderRequest struct {
	Order struct {
		OrganizationID uuid.UUID `json:"-"`
		UserID         uuid.UUID `json:"userID"`
		TableID        uuid.UUID `json:"tableID" binding:"required"`
		UserType       string    `json:"user_type"`
		Notes          string    `json:"notes"`
	}
	OrderItems []struct {
		MenuItemID uuid.UUID `json:"menu_item_id" binding:"required"`
		Quantity   uint64    `json:"quantity" binding:"required"`
	}
}

type GetOrderRequest struct {
	ID             uuid.UUID `json:"id" binding:"required"`
	OrganizationID uuid.UUID `json:"-"`
}

type GetOrganizationOrdersRequest struct {
	OrganizationID uuid.UUID `json:"-"`
}
